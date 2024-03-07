package common

import (
	"context"
	"log"
	"path"
	"reflect"
	"time"

	"github.com/batchatco/go-native-netcdf/netcdf"
	"github.com/batchatco/go-native-netcdf/netcdf/api"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/qiniu/qmgo"
)

type NetCDF struct {
	api      api.Group      `bson:"-"`
	FileName string         `json:"file_name"`
	Data     map[string]any `bson:"data"`
}

func (n NetCDF) GetFloatVar(key string) float64 {
	variable, err := n.api.GetVariable(key)
	if err != nil {
		return 0
	}
	if v, ok := variable.Values.(float64); ok {
		return v
	}
	return 0
}

func (n NetCDF) GetStringSlice(key string) []string {
	variable, err := n.api.GetVariable(key)
	if err != nil {
		return nil
	}
	if v, ok := variable.Values.([]string); ok {
		return v
	}
	return nil
}

func (n NetCDF) GetTime(key string) time.Time {
	variable, err := n.api.GetVariable(key)
	if err != nil {
		log.Println(err.Error())
		return time.Time{}
	}
	if v, ok := variable.Values.(time.Time); ok {
		return v
	}
	return time.Time{}
}

func (n NetCDF) GetGroupFloatSlice(groupName, key string) []float64 {
	group, err := n.api.GetGroup(groupName)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	variable, err := group.GetVariable(key)
	if err != nil {
		log.Println(err.Error())

		return nil
	}

	if v, ok := variable.Values.([]float64); ok {
		return v
	}
	return nil
}

func (n NetCDF) GetGroupTimeSlice(groupName, key string) []time.Time {
	group, err := n.api.GetGroup(groupName)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	variable, err := group.GetVariable(key)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	rv := reflect.ValueOf(variable.Values)
	rvLen := rv.Len()
	if reflect.TypeOf(variable.Values).Kind() == reflect.Slice && rvLen > 0 {
		var timeSlice = make([]time.Time, 0, rvLen)
		switch rv.Index(0).Kind() {
		case reflect.String:
			for i := 0; i < rvLen; i++ {
				s := rv.Index(i).Interface().(string)
				parse, err := time.Parse(time.RFC3339, s)
				if err != nil {
					continue
				}
				timeSlice = append(timeSlice, parse)
			}
			return timeSlice
		case reflect.Int, reflect.Int32, reflect.Int64:
			for i := 0; i < rvLen; i++ {
				nowVal := rv.Index(i).Interface().(int)
				if nowVal > 1000000000000000 {
					timeSlice = append(timeSlice, time.UnixMicro(int64(nowVal)))
				} else if nowVal > 1000000000000 {
					timeSlice = append(timeSlice, time.UnixMilli(int64(nowVal)))
				}
			}
			return timeSlice
		}
	}

	if v, ok := variable.Values.([]time.Time); ok {
		return v
	}
	return nil
}

func (n NetCDF) GetGroupDoubleFloatSlice(groupName, key string) [][]float64 {
	group, err := n.api.GetGroup(groupName)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	variable, err := group.GetVariable(key)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if v, ok := variable.Values.([][]float64); ok {
		return v
	}
	spew.Dump()
	return nil
}

func (n NetCDF) GetGroupDoubleInt32Slice(groupName, key string) [][]int32 {
	group, err := n.api.GetGroup(groupName)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	variable, err := group.GetVariable(key)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	if v, ok := variable.Values.([][]int32); ok {
		return v
	}
	return nil
}

func NewCDF(filepath string) (*NetCDF, error) {
	file, err := netcdf.Open(filepath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var cdf = NetCDF{
		api:  file,
		Data: toJson(file),
	}
	_, cdf.FileName = path.Split(filepath)
	return &cdf, nil
}

func toJson(group api.Group) map[string]any {
	var data = make(map[string]any)
	vars := group.ListVariables()
	var varMap = make(map[string]any, len(vars))
	for _, name := range vars {
		variable, _ := group.GetVariable(name)
		varMap[name] = getVariable(variable)
	}
	data["variable"] = varMap
	groups := group.ListSubgroups()

	var groupMap = make(map[string]any, len(groups))
	for _, groupName := range groups {
		g, _ := group.GetGroup(groupName)
		groupMap[groupName] = toJson(g)
	}
	if len(groupMap) > 0 {
		data["sub_group"] = groupMap
	}

	return data
}

func (n NetCDF) SaveMongo(db *qmgo.Database, tableName string) error {
	_, err := db.Collection(tableName).InsertOne(context.Background(), n.Data)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type Variable struct {
	Value      any            `bson:"value"`
	Attributes map[string]any `bson:"attributes"`
	Dimension  []string       `bson:"dimension"`
}

func getVariable(variable *api.Variable) Variable {
	var val = Variable{
		Value:     variable.Values,
		Dimension: variable.Dimensions,
	}

	keys := variable.Attributes.Keys()
	val.Attributes = make(map[string]any, len(keys))
	for _, key := range keys {
		val.Attributes[key], _ = variable.Attributes.Get(key)
	}
	return val
}
