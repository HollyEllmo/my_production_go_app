package types

// import (
// 	"github.com/HollyEllmo/my-first-go-project/pkg/api/filter"
// 	pb_common_filter "github.com/HollyEllmo/my_proto_repo/gen/go/prod_service/filter/v1"
// )

// type FilterOperator = filter.Operator

// func IntOperatorFromPB(e pb_common_filter.IntFieldFilter_Operator) FilterOperator {
// 	switch e {
// 	case pb_common_filter.IntFieldFilter_OPERATOR_EQ:
// 		return filter.OperatorEq
// 	case pb_common_filter.IntFieldFilter_OPERATOR_NEQ:
// 		return filter.OperatorNotEq
// 	case pb_common_filter.IntFieldFilter_OPERATOR_LT:
// 		return filter.OperatorLowerThan
// 	case pb_common_filter.IntFieldFilter_OPERATOR_LTE:
// 		return filter.OperatorLowerThanEq
// 	case pb_common_filter.IntFieldFilter_OPERATOR_GT:
// 		return filter.OperatorGreaterThan
// 	case pb_common_filter.IntFieldFilter_OPERATOR_GTE:
// 		return filter.OperatorGreaterThanEq
// 	default:
// 		return ""
// 	}
// }