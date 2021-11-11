// ================================================================
// All the usual contents of main() are put into this package for ease of
// testing.
// ================================================================

package entrypoint

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"mlr/internal/pkg/auxents"
	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/climain"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/platform"
	"mlr/internal/pkg/stream"
	"mlr/internal/pkg/transformers"
)

// ----------------------------------------------------------------
func Main() {
	// Special handling for Windows so we can do things like:
	//
	//   mlr put '$a = $b . "cd \"efg\" hi"' foo.dat
	//
	// as on Linux/Unix/MacOS. (On the latter platforms, this is just os.Args
	// as-is.)
	os.Args = platform.GetArgs()

	// Expand "-xyz" into "-x -y -z" while leaving "--xyz" intact. This is a
	// keystroke-saver for the user.
	//
	// This is OK to do globally here since Miller is quite consistent (in
	// main, verbs, and auxents) that multi-character options start with two
	// dashes, e.g. "--csv". (The sole exception is the sort verb's -nf/-nr
	// which are handled specially there.)
	os.Args = lib.Getoptify(os.Args)

	// 'mlr repl' or 'mlr lecat' or any other non-miller-per-se toolery which
	// is delivered (for convenience) within the mlr executable. If argv[1] is
	// found then this function will not return.
	auxents.Dispatch(os.Args)

	options, recordTransformers, err := climain.ParseCommandLine(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)
	}

	if !options.DoInPlace {
		processToStdout(options, recordTransformers)
	} else {
		processInPlace(options)
	}
}

// ----------------------------------------------------------------
// processToStdout is normal processing without mlr -I.

func processToStdout(
	options cli.TOptions,
	recordTransformers []transformers.IRecordTransformer,
) {
	err := stream.Stream(options.FileNames, &options, recordTransformers, os.Stdout, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: %v.\n", err)
		os.Exit(1)
	}
}

// ----------------------------------------------------------------
// processInPlace is in-place processing without mlr -I.
//
// For in-place mode, reconstruct the transformers on each input file. E.g.
// 'mlr -I head -n 2 foo bar' should do head -n 2 on foo as well as on bar.
//
// I could have implemented this with a single construction of the transformers
// and having each transformers implement a Reset() method.  However, having
// effectively two initalizers per transformers -- constructor and reset method
// -- I'd surely miss some logic somewhere.  With in-place mode being a less
// frequently used code path, this would likely lead to latent bugs. So this
// approach leads to greater code stability.

func processInPlace(
	originalOptions cli.TOptions,
) {
	// This should have been already checked by the CLI parser when validating
	// the -I flag.
	lib.InternalCodingErrorIf(originalOptions.FileNames == nil)
	lib.InternalCodingErrorIf(len(originalOptions.FileNames) == 0)

	// Save off the file names from the command line.
	fileNames := make([]string, len(originalOptions.FileNames))
	for i, fileName := range originalOptions.FileNames {
		fileNames[i] = fileName
	}

	for _, fileName := range fileNames {

		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: %v\n", "mlr", err)
			os.Exit(1)
		}

		// Reconstruct the transformers for each file name, and allocate
		// reader, mappers, and writer individually for each file name.  This
		// way CSV headers appear in each file, head -n 10 puts 10 rows for
		// each output file, and so on.
		options, recordTransformers, err := climain.ParseCommandLine(os.Args)
		if err != nil {
			fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
			os.Exit(1)
		}

		// We can't in-place update http://, https://, etc. Also, anything with
		// --prepipe or --prepipex, we won't try to guess how to invert that
		// command to produce re-compressed output.
		err = lib.IsUpdateableInPlace(fileName, options.ReaderOptions.Prepipe)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		containingDirectory := path.Dir(fileName)
		// Names like ./mlr-in-place-2148227797 and ./mlr-in-place-1792078347,
		// as revealed by printing handle.Name().
		handle, err := ioutil.TempFile(containingDirectory, "mlr-in-place-")
		if err != nil {
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		// Run the Miller processing stream from the input file to the temp-output file.
		err = stream.Stream([]string{fileName}, &options, recordTransformers, wrappedHandle, false)
		if err != nil {
			os.Remove(tempFileName)
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		// Close the recompressor handle, if any recompression is being applied.
		if isNew {
			err = wrappedHandle.Close()
			if err != nil {
				os.Remove(tempFileName)
				fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
				os.Exit(1)
			}
		}

		// Close the handle to the output file. This may force final writes, so
		// it must be error-checked.
		err = handle.Close()
		if err != nil {
			os.Remove(tempFileName)
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}

		// Rename the temp-output file on top of the input file.
		err = os.Rename(tempFileName, fileName)
		if err != nil {
			os.Remove(tempFileName)
			fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
			os.Exit(1)
		}
	}
}
