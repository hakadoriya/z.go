// Package sheetz provides a database/sql driver for the Google Sheets API.
//
// Warning: This package is experimental and the syntax is subject to change.
//
// Supported examples:
//
//	SELECT * FROM Sheet1
//	SELECT column1, column2 FROM Sheet1
//	SELECT * FROM Sheet1!A1:Z
//	SELECT column1, column2 FROM "Sheet With Space"!A1:Z
//	(Trailing semicolon is optional.)
package sheetz
