package tdengine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/kubeedge/mapper-generator/pkg/common"
	_ "github.com/taosdata/driver-go/v3/taosRestful"
	"k8s.io/klog/v2"
	"time"
)

var (
	DB *sql.DB
)

type DataBaseConfig struct {
	Config *ConfigData `json:"config,omitempty"`
}
type ConfigData struct {
	Dsn string `json:"dsn,omitempty"`
}

func NewDataBaseClient(config json.RawMessage) (*DataBaseConfig, error) {
	configdata := new(ConfigData)
	err := json.Unmarshal(config, configdata)
	if err != nil {
		return nil, err
	}
	return &DataBaseConfig{
		Config: configdata,
	}, nil
}
func (d *DataBaseConfig) InitDbClient() error {
	var err error
	DB, err = sql.Open("taosRestful", d.Config.Dsn)
	if err != nil {
		//klog.Infof("failed connect to TDengine, err:", err)
		fmt.Printf("failed connect to TDengine:%v", err)
	}
	return nil
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) CloseSessio() {
	err := DB.Close()
	if err != nil {
		klog.Infoln("failded disconnect taosDB")
	}
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) AddData(data *common.DataModel) error {

	stabel := fmt.Sprintf("CREATE STABLE %s (ts timestamp, devicename binary(64), propertyname binary(64), data binary(64),type binary(64)) TAGS (propertyName binary(64));", data.DeviceName)
	fmt.Println(stabel)
	_, err := DB.Exec(stabel)
	if err != nil {
		fmt.Printf("create stable failed %v", err)
	}

	datatime := time.Unix(data.TimeStamp, 0).Format("2006-01-02 15:04:05")
	insertSQL := fmt.Sprintf("INSERT INTO %s USING %s TAGS ('%s') VALUES('%v','%s', '%s', '%s', '%s');",
		data.PropertyName, data.DeviceName, data.PropertyName, datatime, data.DeviceName, data.PropertyName, data.Value, data.Type)

	fmt.Println(insertSQL)
	//tdengine创建超级表第一列必须为时间戳
	_, err = DB.Exec(insertSQL)
	if err != nil {
		klog.Infof("failed add data to tdengine:%v", err)
	}
	return nil
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) GetDataByDeviceName(deviceName string) ([]*common.DataModel, error) {
	query := fmt.Sprintf("SELECT ts, devicename, propertyname, data, type FROM %s", deviceName)
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []*common.DataModel
	for rows.Next() {
		var result common.DataModel
		var ts time.Time
		err := rows.Scan(&ts, &result.DeviceName, &result.PropertyName, &result.Value, &result.Type)
		if err != nil {
			//klog.Infof("scan error:\n", err)
			fmt.Printf("scan error:\n", err)
			return nil, err
		}
		result.TimeStamp = ts.Unix()
		results = append(results, &result)
	}
	return results, nil
	//TODO implement me
	//panic("implement me")
}
func (d *DataBaseConfig) GetPropertyDataByDeviceName(deviceName string, propertyData string) ([]*common.DataModel, error) {
	//TODO implement me
	panic("implement me")
}
func (d *DataBaseConfig) GetDataByTimeRange(start int64, end int64) ([]*common.DataModel, error) {
	//TODO implement me
	panic("implement me")
}
func (d *DataBaseConfig) DeleteDataByTimeRange(start int64, end int64) ([]*common.DataModel, error) {
	//TODO implement me
	panic("implement me")
}
