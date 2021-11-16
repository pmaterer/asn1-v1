# asn1

[ASN.1](https://en.wikipedia.org/wiki/ASN.1) [BER](https://en.wikipedia.org/wiki/X.690#BER_encoding) serialization Go module.

## Basic Encoding Rules (BER)

BER uses a [type-length-value](https://en.wikipedia.org/wiki/Type%E2%80%93length%E2%80%93value) encoding scheme. Encoded data is structured as such:

![](./docs/encoding.png)

The **Identifier octets** identify the *type* of thing encoded. The **Length octets** identify the *length* of the thing encoded. Finally, the **Content octets** contain the *encoded thing*.

## Identifier octets

The identifier octets encode the contents ASN.1 tag's class, whether it is primitive or constructed, and tag number.

![](docs/identifier-octet.png)

### Class

Bits 8 and 7, of the first octet, set the tag's class.

|Class|Bit 8|Bit 7|Description|
|-|-|-|-|
|Universal|0|0|Types native to ASN.1|
|Application|0|1|Types valid for a specific application|
|Context-specific|1|0|Meaning depends on the context (such as within a sequence, set or choice)|
|Private|1|1|Defined in private specifications|

### Primitive/Constructed

Bit 6, of the first octet, sets whether the value is primitive or constructed:

||Bit 6|
|-|-|
|primitive|0|
|constructed|1|

Primitive encodings represent the value directly. For instance, an `INTEGER` tag is primitive, and the underlying content value is an encoded integer.

Constructed encodings represent a concatenation of other encoded values.

### Tag number

For tag numbers <= 30 (the universal tags), the last 5 bits are used to encode the tag number, and the identifier will be 1 octet. For tag numbers >= 31 the last 5 bits of the first octet are encoded to `11111`. The subsequent octets are then the base 128 encoding of the tag number.

### Explicit vs Implicit

Data can be given a unique tag number (especially members of sequences and sets). This helps distinguish that data from other members. These unique tags can be either *explicit* or *implicit*.

#### Explicit

This is the default unique tag style. The data is fully encoded using its underlying type, which is then wrapped in an outer encoding using the unique tag.

## Length octets

BER supports length octets in two possible forms: definite and indefinite. **This module only supports the definite form.**

## Content Encoding

## Types

### Universal class

|ASN.1 Type|Description|Tag|Permitted Construction|Go Types|
|-|-|-|-|-|
|`BOOLEAN`|Simple boolean value|1|Primitive|`bool`|
|`INTEGER`|A signed integer value, with no limits|2|Primitive|`int`, `int8`, `int16`, `int32`, `int64`|
|`BIT STRING`|Arbitrary string of bits|3|Both|`string`|
|`OCTET STRING`|A value of zero or more bytes|4|Both|`[]byte`|
|`NULL`|A empty, non-value|5|Primitive|`asn1.Null`|
|`OBJECT IDENTIFIER`|A sequence of integer components that identify a globally unique object|6|Primitive|`asn1.ObjectIdentifier`|
|`REAL`|A floating point number|9|Primitive|`float32`, `float64`|
|`ENUMERATED`|A numeric value which is associated with a particular meaning|10|Primitive|`asn1.Enumerated`|
|`UTF8String`|UTF-8 string|12|Primitive|`string`|
|`SEQUENCE`|A fixed number of fields of different types, ordered|16|Constructed|`struct`|
|`SEQUENCE OF`|Arbitrary number of fields of different types, ordered|16|Constructed|`[]interface{}`|
|`SET`|A fixed number of fields of different types, unordered|17|Constructed|`struct`|
|`SET OF`|Arbitrary number of fields of different types, unordered|17|Constructed|`[]interface{}`|
|`NumericString`|String representation of a numeric value|18|Primitive|`string`|
|`PrintableString`|A restricted subset of ASCII alphabet|19|Primitive|`string`|
|`IA5String`|First 128 characters of the ASCII alpahabet|22|Primitive|`string`|
|`UTCTime`|Time type, as described [here](https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.5.1)|23|Both|`time.Time`|
|`GeneralizedTime`|Time type, as described [here](https://datatracker.ietf.org/doc/html/rfc5280#section-4.1.2.5.2)|24|Both|`time.Time`|

## Todo

* Handle pointers
* Choice type
* Constructed strings

* Sequence/Set tests

## References

* [A Layman's Guide to a Subset of ASN.1, BER, and DER](https://luca.ntop.org/Teaching/Appunti/asn1.html)
* [LDAPv3 Wire Protocol Reference: The ASN.1 Basic Encoding Rules](https://ldap.com/ldapv3-wire-protocol-reference-asn1-ber/)
* [A Warm Welcome to ASN.1 and DER](https://letsencrypt.org/docs/a-warm-welcome-to-asn1-and-der/#sequence-encoding)