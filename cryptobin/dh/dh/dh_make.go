package dh

import (
    "errors"
)

// 生成公钥
func (this Dh) MakePublicKey() Dh {
    if this.privateKey == nil {
        this.Error = errors.New("Dh: [MakePublicKey()] privateKey error.")
        return this
    }

    // 导出公钥
    this.publicKey = &this.privateKey.PublicKey

    return this
}