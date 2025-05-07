package pkgs

import (
	validator "github.com/go-playground/validator/v10"
)

var DmValidator = validator.New()

const (
	IDEN    = "Identifier"
	CNTRLR  = "controller"
	SVCS    = "services"
	DSTORES = "dmstores"
)
