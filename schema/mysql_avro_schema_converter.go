// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package schema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/uber/storagetapper/db"
	"github.com/uber/storagetapper/types"
)

var (
	namespace = "storagetapper"
)

//MySQLToAvroType is the static conversion map from MySQL types to Avro types
var MySQLToAvroType = map[string]types.AvroPrimitiveType{
	"INT":       types.AvroINT,
	"INTEGER":   types.AvroINT,
	"TINYINT":   types.AvroINT,
	"SMALLINT":  types.AvroINT,
	"MEDIUMINT": types.AvroINT,
	"BOOLEAN":   types.AvroBOOLEAN,
	"BIGINT":    types.AvroLONG,

	"FLOAT":   types.AvroFLOAT,
	"DOUBLE":  types.AvroDOUBLE,
	"DECIMAL": types.AvroDOUBLE,
	"BIT":     types.AvroBYTES,

	"CHAR":    types.AvroSTRING,
	"VARCHAR": types.AvroSTRING,

	//TODO: Confirm correct Avro type for BINARY/VARBINARY types
	"BINARY":    types.AvroBYTES,
	"VARBINARY": types.AvroBYTES,

	"TEXT":       types.AvroSTRING,
	"TINYTEXT":   types.AvroSTRING,
	"MEDIUMTEXT": types.AvroSTRING,
	"LONGTEXT":   types.AvroSTRING,
	"JSON":       types.AvroSTRING,

	"BLOB":       types.AvroBYTES,
	"TINYBLOB":   types.AvroBYTES,
	"MEDIUMBLOB": types.AvroBYTES,
	"LONGBLOB":   types.AvroBYTES,

	"DATE":      types.AvroSTRING,
	"DATETIME":  types.AvroLONG,
	"TIMESTAMP": types.AvroLONG,
	"TIME":      types.AvroSTRING,
	"YEAR":      types.AvroINT,
	"ENUM":      types.AvroSTRING,
	"SET":       types.AvroSTRING,
	//TODO: Add geometry types
}

// ConvertToAvroFromSchema converts a MySQL schema to an Avro schema
func ConvertToAvroFromSchema(tblSchema *types.TableSchema, formatType string) ([]byte, error) {
	avroSchema := &types.AvroSchema{
		Name:      fmt.Sprintf("%s_%s", strings.Replace(tblSchema.DBName, "-", "__", -1), strings.Replace(tblSchema.TableName, "-", "__", -1)),
		Type:      types.AvroRECORD,
		Namespace: namespace,
		Fields:    []types.AvroField{},
		Owner:     tblSchema.DBName,
	}

	for _, colSchema := range tblSchema.Columns {
		avroType := MySQLToAvroType[strings.ToUpper(colSchema.DataType)]
		if avroType == "" {
			continue
		}
		if colSchema.Type == types.MySQLBoolean {
			avroType = MySQLToAvroType["BOOLEAN"]
		}
		fieldTypes := []types.AvroPrimitiveType{types.AvroNULL, avroType}
		avroField := types.AvroField{
			Name:    colSchema.Name,
			Type:    fieldTypes,
			Default: nil,
		}
		avroSchema.Fields = append(avroSchema.Fields, avroField)
	}

	if formatType == "avro" {
		fieldTypes := []types.AvroPrimitiveType{types.AvroLONG}
		avroField := types.AvroField{
			Name:    "ref_key",
			Type:    fieldTypes,
			Default: nil,
		}
		avroSchema.Fields = append(avroSchema.Fields, avroField)

		fieldTypes = []types.AvroPrimitiveType{types.AvroBYTES}
		avroField = types.AvroField{
			Name:    "row_key",
			Type:    fieldTypes,
			Default: nil,
		}
		avroSchema.Fields = append(avroSchema.Fields, avroField)

		fieldTypes = []types.AvroPrimitiveType{types.AvroNULL, types.AvroBOOLEAN}
		avroField = types.AvroField{
			Name:    "is_deleted",
			Type:    fieldTypes,
			Default: nil,
		}
		avroSchema.Fields = append(avroSchema.Fields, avroField)
	}

	return json.Marshal(avroSchema)
}

// ConvertToAvro converts a MySQL schema to an Avro schema
func ConvertToAvro(dbl *db.Loc, tableName string, inputType string, formatType string) ([]byte, error) {
	tblSchema, err := Get(dbl, tableName, inputType)
	if err != nil {
		return nil, err
	}
	return ConvertToAvroFromSchema(tblSchema, formatType)
}
