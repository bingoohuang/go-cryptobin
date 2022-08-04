package pkcs8

import (
    "hash"
    "errors"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "crypto/x509/pkix"
    "encoding/asn1"

    "golang.org/x/crypto/md4"
    "golang.org/x/crypto/pbkdf2"
    "github.com/tjfoc/gmsm/sm3"
)

// pkcs8 可使用的 hash 方式
type Hash uint

const (
    MD2 Hash = 1 + iota // 暂时没有提供
    MD4
    MD5
    SHA1
    SHA224
    SHA256
    SHA384
    SHA512
    SHA512_224
    SHA512_256
    SM3
)

var (
    // key derivation functions
    oidPKCS5          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5}
    oidPKCS5PBKDF2    = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 12}

    // hash 方式
    oidDigestAlgorithm     = asn1.ObjectIdentifier{1, 2, 840, 113549, 2}
    oidHMACWithMD2         = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 2}
    oidHMACWithMD4         = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 4}
    oidHMACWithMD5         = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 5}
    oidHMACWithSHA1        = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 7}
    oidHMACWithSHA224      = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 8}
    oidHMACWithSHA256      = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 9}
    oidHMACWithSHA384      = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 10}
    oidHMACWithSHA512      = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 11}
    oidHMACWithSHA512_224  = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 12}
    oidHMACWithSHA512_256  = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 13}
    oidHMACWithSM3         = asn1.ObjectIdentifier{1, 2, 156, 10197, 1, 401, 2}
)

// 返回使用的 Hash 方式
func prfByOID(oid asn1.ObjectIdentifier) (func() hash.Hash, error) {
    switch {
        case oid.Equal(oidHMACWithMD4):
            return md4.New, nil
        case oid.Equal(oidHMACWithMD5):
            return md5.New, nil
        case oid.Equal(oidHMACWithSHA1):
            return sha1.New, nil
        case oid.Equal(oidHMACWithSHA224):
            return sha256.New224, nil
        case oid.Equal(oidHMACWithSHA256):
            return sha256.New, nil
        case oid.Equal(oidHMACWithSHA384):
            return sha512.New384, nil
        case oid.Equal(oidHMACWithSHA512):
            return sha512.New, nil
        case oid.Equal(oidHMACWithSHA512_224):
            return sha512.New512_224, nil
        case oid.Equal(oidHMACWithSHA512_256):
            return sha512.New512_256, nil
        case oid.Equal(oidHMACWithSM3):
            return sm3.New, nil
    }

    return nil, errors.New("pkcs8: unsupported hash function")
}

// 返回使用的 Hash 对应的 asn1
func oidByHash(h Hash) (asn1.ObjectIdentifier, error) {
    switch h {
        case MD4:
            return oidHMACWithMD4, nil
        case MD5:
            return oidHMACWithMD5, nil
        case SHA1:
            return oidHMACWithSHA1, nil
        case SHA224:
            return oidHMACWithSHA224, nil
        case SHA256:
            return oidHMACWithSHA256, nil
        case SHA384:
            return oidHMACWithSHA384, nil
        case SHA512:
            return oidHMACWithSHA512, nil
        case SHA512_224:
            return oidHMACWithSHA512_224, nil
        case SHA512_256:
            return oidHMACWithSHA512_256, nil
        case SM3:
            return oidHMACWithSM3, nil
    }

    return nil, errors.New("pkcs8: unsupported hash function")
}

// pbkdf2 数据
type pbkdf2Params struct {
    Salt           []byte
    IterationCount int
    PrfParam       pkix.AlgorithmIdentifier `asn1:"optional"`
}

func (this pbkdf2Params) DeriveKey(password []byte, size int) (key []byte, err error) {
    h, err := prfByOID(this.PrfParam.Algorithm)
    if err != nil {
        return nil, err
    }

    return pbkdf2.Key(password, this.Salt, this.IterationCount, size, h), nil
}

// PBKDF2 配置
type PBKDF2Opts struct {
    SaltSize       int
    IterationCount int
    HMACHash       Hash
}

func (this PBKDF2Opts) DeriveKey(password, salt []byte, size int) (key []byte, params KDFParameters, err error) {
    alg, err := oidByHash(this.HMACHash)
    if err != nil {
        return nil, nil, err
    }

    h, err := prfByOID(alg)
    if err != nil {
        return nil, nil, err
    }

    key = pbkdf2.Key(password, salt, this.IterationCount, size, h)

    prfParam := pkix.AlgorithmIdentifier{
        Algorithm:  alg,
        Parameters: asn1.RawValue{
            Tag: asn1.TagNull,
        },
    }

    params = pbkdf2Params{salt, this.IterationCount, prfParam}
    return key, params, nil
}

func (this PBKDF2Opts) GetSaltSize() int {
    return this.SaltSize
}

func (this PBKDF2Opts) OID() asn1.ObjectIdentifier {
    return oidPKCS5PBKDF2
}

func init() {
    AddKDF(oidPKCS5PBKDF2, new(pbkdf2Params))
}
