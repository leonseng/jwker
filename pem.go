package main

import (
  "crypto/ecdsa"
  "crypto/rsa"
  "crypto/x509"
  "encoding/pem"
)

func pemToJwk(pemBytes []byte) string {
  pemBlock := findPemBlock(pemBytes)
  keyStruct := processBlock(pemBlock)
  return structToJWK(keyStruct)
}

func findPemBlock(pemBytes []byte) *pem.Block {
  // Only get the first PEM block
  pemBlock, _ := pem.Decode(pemBytes)

  if pemBlock == nil {
    throwParseError("invalid PEM file format.")
  }

  if (x509.IsEncryptedPEMBlock(pemBlock)) {
    throwParseError("the given PEM file is encrypted. Please decrypt first.")
  }

  return pemBlock
}

func processBlock(pemBlock *pem.Block) interface{} {
  var keyStruct interface{}

  switch pemBlock.Type  {
  case "PUBLIC KEY":
    keyStruct = processPublicKey(pemBlock.Bytes)
  case "RSA PRIVATE KEY":
    key, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
    stopOnParseError(err)
    keyStruct = processRSAPrivate(key)
  default:
    throwParseError("unsupported PEM type.")
  }

  return keyStruct
}

func processPublicKey(bytes []byte) interface{} {
  key, err := x509.ParsePKIXPublicKey(bytes)
  stopOnParseError(err)

  var keyStruct interface{}

  switch key := key.(type) {
  case *rsa.PublicKey:
    keyStruct = processRSAPublic(key)
  case *ecdsa.PublicKey:
    keyStruct = processECPublic(key)
  default:
    throwParseError("Unknown key type.")
  }

  return keyStruct
}