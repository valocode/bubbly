package adapter

import (
	"os"

	"github.com/hashicorp/hcl/v2/ext/tryfunc"
	ctyyaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"

	"github.com/hashicorp/terraform/lang/funcs"
)

// stdfunctions returns functions for the SymbolTable's EvalContext
func stdfunctions() map[string]function.Function {
	return map[string]function.Function{
		// Our own custom functions here
		"env": EnvFunc,

		// The following are from cty stdlib that we pull in
		"and": stdlib.AndFunc,
		"or":  stdlib.OrFunc,

		//
		// The rest are straight from Terraform:
		// https://github.com/hashicorp/terraform/blob/main/lang/functions.go
		//
		"abs":          stdlib.AbsoluteFunc,
		"abspath":      funcs.AbsPathFunc,
		"alltrue":      funcs.AllTrueFunc,
		"anytrue":      funcs.AnyTrueFunc,
		"basename":     funcs.BasenameFunc,
		"base64decode": funcs.Base64DecodeFunc,
		"base64encode": funcs.Base64EncodeFunc,
		"base64gzip":   funcs.Base64GzipFunc,
		"base64sha256": funcs.Base64Sha256Func,
		"base64sha512": funcs.Base64Sha512Func,
		"bcrypt":       funcs.BcryptFunc,
		"can":          tryfunc.CanFunc,
		"ceil":         stdlib.CeilFunc,
		"chomp":        stdlib.ChompFunc,
		"cidrhost":     funcs.CidrHostFunc,
		"cidrnetmask":  funcs.CidrNetmaskFunc,
		"cidrsubnet":   funcs.CidrSubnetFunc,
		"cidrsubnets":  funcs.CidrSubnetsFunc,
		"coalesce":     funcs.CoalesceFunc,
		"coalescelist": stdlib.CoalesceListFunc,
		"compact":      stdlib.CompactFunc,
		"concat":       stdlib.ConcatFunc,
		"contains":     stdlib.ContainsFunc,
		"csvdecode":    stdlib.CSVDecodeFunc,
		// "defaults":         s.experimentalFunction(experiments.ModuleVariableOptionalAttrs, funcs.DefaultsFunc),
		"dirname":          funcs.DirnameFunc,
		"distinct":         stdlib.DistinctFunc,
		"element":          stdlib.ElementFunc,
		"chunklist":        stdlib.ChunklistFunc,
		"file":             funcs.MakeFileFunc(".", false),
		"fileexists":       funcs.MakeFileExistsFunc("."),
		"fileset":          funcs.MakeFileSetFunc("."),
		"filebase64":       funcs.MakeFileFunc(".", true),
		"filebase64sha256": funcs.MakeFileBase64Sha256Func("."),
		"filebase64sha512": funcs.MakeFileBase64Sha512Func("."),
		"filemd5":          funcs.MakeFileMd5Func("."),
		"filesha1":         funcs.MakeFileSha1Func("."),
		"filesha256":       funcs.MakeFileSha256Func("."),
		"filesha512":       funcs.MakeFileSha512Func("."),
		"flatten":          stdlib.FlattenFunc,
		"floor":            stdlib.FloorFunc,
		"format":           stdlib.FormatFunc,
		"formatdate":       stdlib.FormatDateFunc,
		"formatlist":       stdlib.FormatListFunc,
		"indent":           stdlib.IndentFunc,
		"index":            funcs.IndexFunc, // stdlib.IndexFunc is not compatible
		"join":             stdlib.JoinFunc,
		"jsondecode":       stdlib.JSONDecodeFunc,
		"jsonencode":       stdlib.JSONEncodeFunc,
		"keys":             stdlib.KeysFunc,
		"length":           funcs.LengthFunc,
		"list":             funcs.ListFunc,
		"log":              stdlib.LogFunc,
		"lookup":           funcs.LookupFunc,
		"lower":            stdlib.LowerFunc,
		"map":              funcs.MapFunc,
		"matchkeys":        funcs.MatchkeysFunc,
		"max":              stdlib.MaxFunc,
		"md5":              funcs.Md5Func,
		"merge":            stdlib.MergeFunc,
		"min":              stdlib.MinFunc,
		"one":              funcs.OneFunc,
		"parseint":         stdlib.ParseIntFunc,
		"pathexpand":       funcs.PathExpandFunc,
		"pow":              stdlib.PowFunc,
		"range":            stdlib.RangeFunc,
		"regex":            stdlib.RegexFunc,
		"regexall":         stdlib.RegexAllFunc,
		"replace":          funcs.ReplaceFunc,
		"reverse":          stdlib.ReverseListFunc,
		"rsadecrypt":       funcs.RsaDecryptFunc,
		"sensitive":        funcs.SensitiveFunc,
		"nonsensitive":     funcs.NonsensitiveFunc,
		"setintersection":  stdlib.SetIntersectionFunc,
		"setproduct":       stdlib.SetProductFunc,
		"setsubtract":      stdlib.SetSubtractFunc,
		"setunion":         stdlib.SetUnionFunc,
		"sha1":             funcs.Sha1Func,
		"sha256":           funcs.Sha256Func,
		"sha512":           funcs.Sha512Func,
		"signum":           stdlib.SignumFunc,
		"slice":            stdlib.SliceFunc,
		"sort":             stdlib.SortFunc,
		"split":            stdlib.SplitFunc,
		"strrev":           stdlib.ReverseFunc,
		"substr":           stdlib.SubstrFunc,
		"sum":              funcs.SumFunc,
		"textdecodebase64": funcs.TextDecodeBase64Func,
		"textencodebase64": funcs.TextEncodeBase64Func,
		"timestamp":        funcs.TimestampFunc,
		"timeadd":          stdlib.TimeAddFunc,
		"title":            stdlib.TitleFunc,
		"tostring":         funcs.MakeToFunc(cty.String),
		"tonumber":         funcs.MakeToFunc(cty.Number),
		"tobool":           funcs.MakeToFunc(cty.Bool),
		"toset":            funcs.MakeToFunc(cty.Set(cty.DynamicPseudoType)),
		"tolist":           funcs.MakeToFunc(cty.List(cty.DynamicPseudoType)),
		"tomap":            funcs.MakeToFunc(cty.Map(cty.DynamicPseudoType)),
		"transpose":        funcs.TransposeFunc,
		"trim":             stdlib.TrimFunc,
		"trimprefix":       stdlib.TrimPrefixFunc,
		"trimspace":        stdlib.TrimSpaceFunc,
		"trimsuffix":       stdlib.TrimSuffixFunc,
		"try":              tryfunc.TryFunc,
		"upper":            stdlib.UpperFunc,
		"urlencode":        funcs.URLEncodeFunc,
		"uuid":             funcs.UUIDFunc,
		"uuidv5":           funcs.UUIDV5Func,
		"values":           stdlib.ValuesFunc,
		"yamldecode":       ctyyaml.YAMLDecodeFunc,
		"yamlencode":       ctyyaml.YAMLEncodeFunc,
		"zipmap":           stdlib.ZipmapFunc,
	}
}

// EnvFunc is a function that looks up the value of an environment variable.
// If the variable is undefined, an empty string is returned instead.
var EnvFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "name",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (ret cty.Value, err error) {
		return cty.StringVal(os.Getenv(args[0].AsString())), nil
	},
})
