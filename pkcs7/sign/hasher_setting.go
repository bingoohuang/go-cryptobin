package sign

import (
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/asn1"
)

var (
    // Digest Algorithms
    oidDigestAlgorithmMd5    = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 5}
    oidDigestAlgorithmSHA1   = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 26}
    oidDigestAlgorithmSHA256 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}
    oidDigestAlgorithmSHA384 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 2}
    oidDigestAlgorithmSHA512 = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 3}
)

// 各种 hash
var SignHashWithSHA1 = SignHashWithFunc{
    hashFunc:   sha1.New,
    identifier: oidDigestAlgorithmSHA1,
}

var SignHashWithSHA256 = SignHashWithFunc{
    hashFunc:   sha256.New,
    identifier: oidDigestAlgorithmSHA256,
}

var SignHashWithSHA384 = SignHashWithFunc{
    hashFunc:   sha512.New384,
    identifier: oidDigestAlgorithmSHA384,
}

var SignHashWithSHA512 = SignHashWithFunc{
    hashFunc:   sha512.New,
    identifier: oidDigestAlgorithmSHA512,
}

func init() {
    AddSignHash(oidDigestAlgorithmSHA1, func() SignHash {
        return SignHashWithSHA1
    })
    AddSignHash(oidDigestAlgorithmSHA256, func() SignHash {
        return SignHashWithSHA256
    })
    AddSignHash(oidDigestAlgorithmSHA384, func() SignHash {
        return SignHashWithSHA384
    })
    AddSignHash(oidDigestAlgorithmSHA512, func() SignHash {
        return SignHashWithSHA512
    })
}
