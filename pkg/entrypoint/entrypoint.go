// All the usual contents of main() are put into this package for ease of
// testing.

package entrypoint

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/auxents"
	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/climain"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/platform"
	"github.com/johnkerl/miller/v6/pkg/stream"
	"github.com/johnkerl/miller/v6/pkg/transformers"
)

type MainReturn struct {
	PrintElapsedTime bool
}

func Main() MainReturn {
	// Special handling for Windows so we can do things like:
	//
	//   mlr put '$a = $b . "cd \"efg\" hi"' foo.dat
	//
	// as on Linux/Unix/MacOS. (On the latter platforms, this is just os.Args
	// as-is.)
	os.Args = platform.GetArgs()

	// Enable ANSI escape sequence processing on the Windows pseudo terminal,
	// otherwise, we only raw ANSI escape sequences like ←[0;30m  0←[0m ←[0;31m  1
	platform.EnableAnsiEscapeSequences()

	// 'mlr repl' or 'mlr lecat' or any other non-miller-per-se toolery which
	// is delivered (for convenience) within the mlr executable. If argv[1] is
	// found then this function will not return.
	auxents.Dispatch(os.Args)

	options, recordTransformers, err := climain.ParseCommandLine(os.Args)
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	if !options.DoInPlace {
		err = processToStdout(options, recordTransformers)
	} else {
		err = processFilesInPlace(options)
	}
	if err != nil {
		printError(err)
		os.Exit(1)
	}

	return MainReturn{
		PrintElapsedTime: options.PrintElapsedTime,
	}
}

// printError prints err to stderr. Errors that already start with "mlr " (e.g.
// "mlr stats1: ...") are printed as-is to avoid double-prefixing.
func printError(err error) {
	if strings.HasPrefix(err.Error(), "mlr ") {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
	}
}

// processToStdout is normal processing without mlr -I.

func processToStdout(
	options *cli.TOptions,
	recordTransformers []transformers.RecordTransformer,
) error {
	return stream.Stream(options.FileNames, options, recordTransformers, os.Stdout, true)
}

// processFilesInPlace is in-place processing without mlr -I.
//
// For in-place mode, reconstruct the transformers on each input file. E.g.
// 'mlr -I head -n 2 foo bar' should do head -n 2 on foo as well as on bar.
//
// I could have implemented this with a single construction of the transformers
// and having each transformers implement a Reset() method.  However, having
// effectively two initializers per transformers -- constructor and reset method
// -- I'd surely miss some logic somewhere.  With in-place mode being a less
// frequently used code path, this would likely lead to latent bugs. So this
// approach leads to greater code stability.

func processFilesInPlace(
	originalOptions *cli.TOptions,
) error {
	// This should have been already checked by the CLI parser when validating
	// the -I flag.
	lib.InternalCodingErrorIf(originalOptions.FileNames == nil)
	lib.InternalCodingErrorIf(len(originalOptions.FileNames) == 0)

	// Save off the file names from the command line.
	fileNames := make([]string, len(originalOptions.FileNames))
	copy(fileNames, originalOptions.FileNames)

	for _, fileName := range fileNames {
		err := processFileInPlace(fileName, originalOptions)
		if err != nil {
			return err
		}
	}
	return nil
}

func processFileInPlace(
	fileName string,
	originalOptions *cli.TOptions,
) error {

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return err
	}

	// Reconstruct the transformers for each file name, and allocate
	// reader, mappers, and writer individually for each file name.  This
	// way CSV headers appear in each file, head -n 10 puts 10 rows for
	// each output file, and so on.
	options, recordTransformers, err := climain.ParseCommandLine(os.Args)
	if err != nil {
		return err
	}

	// We can't in-place update http://, https://, etc. Also, anything with
	// --prepipe or --prepipex, we won't try to guess how to invert that
	// command to produce re-compressed output.
	err = lib.IsUpdateableInPlace(fileName, options.ReaderOptions.Prepipe)
	if err != nil {
		return err
	}

	// Get the original file's mode so we can preserve it.
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	originalMode := fileInfo.Mode()

	containingDirectory := path.Dir(fileName)
	// Names like ./mlr-in-place-2148227797 and ./mlr-in-place-1792078347,
	// as revealed by printing handle.Name().
	handle, err := os.CreateTemp(containingDirectory, "mlr-in-place-")
	if err != nil {
		return err
	}
	tempFileName := handle.Name()

	// If the input file is compressed and we'll be doing in-process
	// decompression as we read the input file, try to do in-process
	// compression as we write the output.
	inputFileEncoding := lib.FindInputEncoding(fileName, options.ReaderOptions.FileInputEncoding)

	// Get a handle with, perhaps, a recompression wrapper around it.
	wrappedHandle, isNew, err := lib.WrapOutputHandle(handle, inputFileEncoding)
	if err != nil {
		os.Remove(tempFileName)
		return err
	}

	// Run the Miller processing stream from the input file to the temp-output file.
	err = stream.Stream([]string{fileName}, options, recordTransformers, wrappedHandle, false)
	if err != nil {
		os.Remove(tempFileName)
		return err
	}

	// Close the recompressor handle, if any recompression is being applied.
	if isNew {
		err = wrappedHandle.Close()
		if err != nil {
			os.Remove(tempFileName)
			return err
		}
	}

	// Close the handle to the output file. This may force final writes, so
	// it must be error-checked.
	err = handle.Close()
	if err != nil {
		os.Remove(tempFileName)
		return err
	}

	// Rename the temp-output file on top of the input file.
	err = os.Rename(tempFileName, fileName)
	if err != nil {
		os.Remove(tempFileName)
		return err
	}

	// Set the mode to match the original.
	err = os.Chmod(fileName, originalMode)
	if err != nil {
		return err
	}

	return nil
}
