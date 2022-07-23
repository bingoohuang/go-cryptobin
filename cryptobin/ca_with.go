package cryptobin

import (
    "crypto/x509"
)

// 设置 csr
func (this CA) WithCsr(data *x509.Certificate) CA {
    this.csr = data

    return this
}

// 设置 PrivateKey
func (this CA) WithPrivateKey(data any) CA {
    this.privateKey = data

    return this
}

// 设置 publicKey
func (this CA) WithPublicKey(data any) CA {
    this.publicKey = data

    return this
}

// 设置 keyData
func (this CA) WithKeyData(data []byte) CA {
    this.keyData = data

    return this
}

// 设置错误
func (this CA) WithError(err error) CA {
    this.Error = err

    return this
}