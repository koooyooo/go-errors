package errors

// Domains
var LabelDomainAuthorization = Label("Authorization")
var LabelDomainLogging = Label("Logging")

// Operation Type
var LabelWithGlobalCache = Label("w/GlobalCache")
var LabelWithLock = Label("w/Lock")
var LabelWithGoroutine = Label("w/Goroutine")
var LabelWithFileAccess = Label("w/FileAccess")
var LabelWithJsonConversion = Label("w/Json-Conversion")

var LabelWithAPICallInternal = Label("w/API-Call-Internal")
var LabelWithAPICallExternal = Label("w/API-Call-External")

