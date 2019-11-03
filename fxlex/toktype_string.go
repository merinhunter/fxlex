// Code generated by "stringer -type TokType fxlex.go"; DO NOT EDIT.

package fxlex

import "strconv"

const (
	_TokType_name_0 = "TokKeyTokIdTokFuncTokIntTokIntLitTokBoolTokBoolLitTokCoord"
	_TokType_name_1 = "Declaration"
	_TokType_name_2 = "TokPowTokGTETokLTE"
	_TokType_name_3 = "TokNeg"
	_TokType_name_4 = "TokRemTokAnd"
	_TokType_name_5 = "TokLParTokRParTokTimesTokPlusTokCommaTokMinusTokDotTokDivide"
	_TokType_name_6 = "SemicolonTokLTAssignationTokGT"
	_TokType_name_7 = "TokLSquare"
	_TokType_name_8 = "TokRSquareTokXor"
	_TokType_name_9 = "TokLCurlTokOrTokRCurl"
)

var (
	_TokType_index_0 = [...]uint8{0, 6, 11, 18, 24, 33, 40, 50, 58}
	_TokType_index_2 = [...]uint8{0, 6, 12, 18}
	_TokType_index_4 = [...]uint8{0, 6, 12}
	_TokType_index_5 = [...]uint8{0, 7, 14, 22, 29, 37, 45, 51, 60}
	_TokType_index_6 = [...]uint8{0, 9, 14, 25, 30}
	_TokType_index_8 = [...]uint8{0, 10, 16}
	_TokType_index_9 = [...]uint8{0, 8, 13, 21}
)

func (i TokType) String() string {
	switch {
	case 1 <= i && i <= 8:
		i -= 1
		return _TokType_name_0[_TokType_index_0[i]:_TokType_index_0[i+1]]
	case i == 19:
		return _TokType_name_1
	case 27 <= i && i <= 29:
		i -= 27
		return _TokType_name_2[_TokType_index_2[i]:_TokType_index_2[i+1]]
	case i == 33:
		return _TokType_name_3
	case 37 <= i && i <= 38:
		i -= 37
		return _TokType_name_4[_TokType_index_4[i]:_TokType_index_4[i+1]]
	case 40 <= i && i <= 47:
		i -= 40
		return _TokType_name_5[_TokType_index_5[i]:_TokType_index_5[i+1]]
	case 59 <= i && i <= 62:
		i -= 59
		return _TokType_name_6[_TokType_index_6[i]:_TokType_index_6[i+1]]
	case i == 91:
		return _TokType_name_7
	case 93 <= i && i <= 94:
		i -= 93
		return _TokType_name_8[_TokType_index_8[i]:_TokType_index_8[i+1]]
	case 123 <= i && i <= 125:
		i -= 123
		return _TokType_name_9[_TokType_index_9[i]:_TokType_index_9[i+1]]
	default:
		return "TokType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}