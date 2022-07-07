package version

// STRING is the current Miller major/minor/patch version as a single string.
// Nominally things like "6.0.0" for a release, then "6.0.0-dev" in between.
// This makes it clear that a given build is on the main dev branch, not a
// particular snapshot tag.
var STRING string = "6.3.0"
