// This is a generated source file. DO NOT EDIT
// Source: strfmt/strfmt_generated.go

package strfmt

import "github.com/machinefi/w3bstream/pkg/depends/kit/validator"

func init() {
	validator.DefaultFactory.Register(ASCIIValidator)
}

var ASCIIValidator = validator.NewRegexpStrfmtValidator(regexpStringASCII, "ascii")

func init() {
	validator.DefaultFactory.Register(AlphaValidator)
}

var AlphaValidator = validator.NewRegexpStrfmtValidator(regexpStringAlpha, "alpha")

func init() {
	validator.DefaultFactory.Register(AlphaNumericValidator)
}

var AlphaNumericValidator = validator.NewRegexpStrfmtValidator(regexpStringAlphaNumeric, "alpha-numeric", "alphaNumeric")

func init() {
	validator.DefaultFactory.Register(AlphaUnicodeValidator)
}

var AlphaUnicodeValidator = validator.NewRegexpStrfmtValidator(regexpStringAlphaUnicode, "alpha-unicode", "alphaUnicode")

func init() {
	validator.DefaultFactory.Register(AlphaUnicodeNumericValidator)
}

var AlphaUnicodeNumericValidator = validator.NewRegexpStrfmtValidator(regexpStringAlphaUnicodeNumeric, "alpha-unicode-numeric", "alphaUnicodeNumeric")

func init() {
	validator.DefaultFactory.Register(Base64Validator)
}

var Base64Validator = validator.NewRegexpStrfmtValidator(regexpStringBase64, "base-64", "base64")

func init() {
	validator.DefaultFactory.Register(Base64URLValidator)
}

var Base64URLValidator = validator.NewRegexpStrfmtValidator(regexpStringBase64URL, "base-64-url", "base64Url")

func init() {
	validator.DefaultFactory.Register(BtcAddressValidator)
}

var BtcAddressValidator = validator.NewRegexpStrfmtValidator(regexpStringBtcAddress, "btc-address", "btcAddress")

func init() {
	validator.DefaultFactory.Register(BtcAddressLowerBech32Validator)
}

var BtcAddressLowerBech32Validator = validator.NewRegexpStrfmtValidator(regexpStringBtcAddressLowerBech32, "btc-address-lower-bech-32", "btcAddressLowerBech32")

func init() {
	validator.DefaultFactory.Register(BtcAddressUpperBech32Validator)
}

var BtcAddressUpperBech32Validator = validator.NewRegexpStrfmtValidator(regexpStringBtcAddressUpperBech32, "btc-address-upper-bech-32", "btcAddressUpperBech32")

func init() {
	validator.DefaultFactory.Register(DataURIValidator)
}

var DataURIValidator = validator.NewRegexpStrfmtValidator(regexpStringDataURI, "data-uri", "dataUri")

func init() {
	validator.DefaultFactory.Register(E164Validator)
}

var E164Validator = validator.NewRegexpStrfmtValidator(regexpStringE164, "e-164", "e164")

func init() {
	validator.DefaultFactory.Register(EmailValidator)
}

var EmailValidator = validator.NewRegexpStrfmtValidator(regexpStringEmail, "email")

func init() {
	validator.DefaultFactory.Register(EthAddressValidator)
}

var EthAddressValidator = validator.NewRegexpStrfmtValidator(regexpStringEthAddress, "eth-address", "ethAddress")

func init() {
	validator.DefaultFactory.Register(EthAddressLowerValidator)
}

var EthAddressLowerValidator = validator.NewRegexpStrfmtValidator(regexpStringEthAddressLower, "eth-address-lower", "ethAddressLower")

func init() {
	validator.DefaultFactory.Register(EthAddressUpperValidator)
}

var EthAddressUpperValidator = validator.NewRegexpStrfmtValidator(regexpStringEthAddressUpper, "eth-address-upper", "ethAddressUpper")

func init() {
	validator.DefaultFactory.Register(HslValidator)
}

var HslValidator = validator.NewRegexpStrfmtValidator(regexpStringHSL, "hsl")

func init() {
	validator.DefaultFactory.Register(HslaValidator)
}

var HslaValidator = validator.NewRegexpStrfmtValidator(regexpStringHSLA, "hsla")

func init() {
	validator.DefaultFactory.Register(HTMLValidator)
}

var HTMLValidator = validator.NewRegexpStrfmtValidator(regexpStringHTML, "html")

func init() {
	validator.DefaultFactory.Register(HTMLEncodedValidator)
}

var HTMLEncodedValidator = validator.NewRegexpStrfmtValidator(regexpStringHTMLEncoded, "html-encoded", "htmlEncoded")

func init() {
	validator.DefaultFactory.Register(HexColorValidator)
}

var HexColorValidator = validator.NewRegexpStrfmtValidator(regexpStringHexColor, "hex-color", "hexColor")

func init() {
	validator.DefaultFactory.Register(HexadecimalValidator)
}

var HexadecimalValidator = validator.NewRegexpStrfmtValidator(regexpStringHexadecimal, "hexadecimal")

func init() {
	validator.DefaultFactory.Register(HostnameRfc1123Validator)
}

var HostnameRfc1123Validator = validator.NewRegexpStrfmtValidator(regexpStringHostnameRFC1123, "hostname-rfc-1123", "hostnameRfc1123")

func init() {
	validator.DefaultFactory.Register(HostnameRfc952Validator)
}

var HostnameRfc952Validator = validator.NewRegexpStrfmtValidator(regexpStringHostnameRFC952, "hostname-rfc-952", "hostnameRfc952")

func init() {
	validator.DefaultFactory.Register(Isbn10Validator)
}

var Isbn10Validator = validator.NewRegexpStrfmtValidator(regexpStringISBN10, "isbn-10", "isbn10")

func init() {
	validator.DefaultFactory.Register(Isbn13Validator)
}

var Isbn13Validator = validator.NewRegexpStrfmtValidator(regexpStringISBN13, "isbn-13", "isbn13")

func init() {
	validator.DefaultFactory.Register(LatitudeValidator)
}

var LatitudeValidator = validator.NewRegexpStrfmtValidator(regexpStringLatitude, "latitude")

func init() {
	validator.DefaultFactory.Register(LongitudeValidator)
}

var LongitudeValidator = validator.NewRegexpStrfmtValidator(regexpStringLongitude, "longitude")

func init() {
	validator.DefaultFactory.Register(MultibyteValidator)
}

var MultibyteValidator = validator.NewRegexpStrfmtValidator(regexpStringMultibyte, "multibyte")

func init() {
	validator.DefaultFactory.Register(NumberValidator)
}

var NumberValidator = validator.NewRegexpStrfmtValidator(regexpStringNumber, "number")

func init() {
	validator.DefaultFactory.Register(NumericValidator)
}

var NumericValidator = validator.NewRegexpStrfmtValidator(regexpStringNumeric, "numeric")

func init() {
	validator.DefaultFactory.Register(PrintableASCIIValidator)
}

var PrintableASCIIValidator = validator.NewRegexpStrfmtValidator(regexpStringPrintableASCII, "printable-ascii", "printableAscii")

func init() {
	validator.DefaultFactory.Register(RgbValidator)
}

var RgbValidator = validator.NewRegexpStrfmtValidator(regexpStringRGB, "rgb")

func init() {
	validator.DefaultFactory.Register(RgbaValidator)
}

var RgbaValidator = validator.NewRegexpStrfmtValidator(regexpStringRGBA, "rgba")

func init() {
	validator.DefaultFactory.Register(SsnValidator)
}

var SsnValidator = validator.NewRegexpStrfmtValidator(regexpStringSSN, "ssn")

func init() {
	validator.DefaultFactory.Register(URLValidator)
}

var URLValidator = validator.NewRegexpStrfmtValidator(regexpStringURL, "url")

func init() {
	validator.DefaultFactory.Register(URLEncodedValidator)
}

var URLEncodedValidator = validator.NewRegexpStrfmtValidator(regexpStringURLEncoded, "url-encoded", "urlEncoded")

func init() {
	validator.DefaultFactory.Register(UUIDValidator)
}

var UUIDValidator = validator.NewRegexpStrfmtValidator(regexpStringUUID, "uuid")

func init() {
	validator.DefaultFactory.Register(UUID3Validator)
}

var UUID3Validator = validator.NewRegexpStrfmtValidator(regexpStringUUID3, "uuid-3", "uuid3")

func init() {
	validator.DefaultFactory.Register(UUID3Rfc4122Validator)
}

var UUID3Rfc4122Validator = validator.NewRegexpStrfmtValidator(regexpStringUUID3RFC4122, "uuid-3-rfc-4122", "uuid3Rfc4122")

func init() {
	validator.DefaultFactory.Register(UUID4Validator)
}

var UUID4Validator = validator.NewRegexpStrfmtValidator(regexpStringUUID4, "uuid-4", "uuid4")

func init() {
	validator.DefaultFactory.Register(UUID4Rfc4122Validator)
}

var UUID4Rfc4122Validator = validator.NewRegexpStrfmtValidator(regexpStringUUID4RFC4122, "uuid-4-rfc-4122", "uuid4Rfc4122")

func init() {
	validator.DefaultFactory.Register(UUID5Validator)
}

var UUID5Validator = validator.NewRegexpStrfmtValidator(regexpStringUUID5, "uuid-5", "uuid5")

func init() {
	validator.DefaultFactory.Register(UUID5Rfc4122Validator)
}

var UUID5Rfc4122Validator = validator.NewRegexpStrfmtValidator(regexpStringUUID5RFC4122, "uuid-5-rfc-4122", "uuid5Rfc4122")

func init() {
	validator.DefaultFactory.Register(Uuidrfc4122Validator)
}

var Uuidrfc4122Validator = validator.NewRegexpStrfmtValidator(regexpStringUUIDRFC4122, "uuidrfc-4122", "uuidrfc4122")
