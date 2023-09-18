package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/kubeedge/mapper-generator/pkg/common"
	"k8s.io/klog/v2"
	"strconv"
)

var (
	RedisCli *redis.Client
)

type DataBaseConfig struct {
	Config *ConfigData
}

type ConfigData struct {
	Addr         string `json:"addr,omitempty"`
	Password     string `json:"password,omitempty"`
	DB           int64  `json:"db,omitempty"`
	PoolSize     int64  `json:"poolSize,omitempty"`
	MinIdleConns int64  `json:"minIdleConns,omitempty"`
}

func NewDataBaseClient(config json.RawMessage) (*DataBaseConfig, error) {
	configdata := new(ConfigData)
	err := json.Unmarshal(config, configdata)
	if err != nil {
		return nil, err
	}
	return &DataBaseConfig{Config: configdata}, nil
}

func (d *DataBaseConfig) InitDbClient() error {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:         d.Config.Addr,
		Password:     d.Config.Password,
		DB:           int(d.Config.DB),
		PoolSize:     int(d.Config.PoolSize),
		MinIdleConns: int(d.Config.MinIdleConns),
	})
	pong, err := RedisCli.Ping(context.Background()).Result()
	if err != nil {
		klog.Errorf("init redis database failed, err = %v", err)
		return err
	} else {
		klog.V(1).Infof("init redis database successfully, with return cmd %s", pong)
	}
	return nil
}

func (d *DataBaseConfig) CloseSession() {
	err := RedisCli.Close()
	if err != nil {
		klog.V(4).Info("close database failed")
	}
}

func (d *DataBaseConfig) AddData(data *common.DataModel) error {
	ctx := context.Background()
	// The key to construct the ordered set, here DeviceName is used as the key
	klog.V(1).Infof("deviceName:%s", data.DeviceName)
	// Check if the current ordered set exists
	exists, err := RedisCli.Exists(ctx, data.DeviceName).Result()
	if err != nil {
		klog.V(4).Info("Exit AddData")
		return err
	}
	deviceData := "TimeStamp: " + strconv.FormatInt(data.TimeStamp, 10) + " PropertyName: " + data.PropertyName + " data: " + data.Value
	if exists == 0 {
		// The ordered set does not exist, create a new ordered set and add data
		_, err = RedisCli.ZAdd(ctx, data.DeviceName, &redis.Z{
			Score:  float64(data.TimeStamp),
			Member: deviceData,
		}).Result()
		if err != nil {
			klog.V(4).Info("Exit AddData")
			return err
		}
	} else {
		// The ordered set already exists, add data directly
		_, err = RedisCli.ZAdd(ctx, data.DeviceName, &redis.Z{
			Score:  float64(data.TimeStamp),
			Member: deviceData,
		}).Result()
		if err != nil {
			klog.V(4).Info("Exit AddData")
			return err
		}
	}
	return nil
}

func (d *DataBaseConfig) GetDataByDeviceName(deviceName string) ([]*common.DataModel, error) {
	ctx := context.Background()

	dataJSON, err := RedisCli.ZRevRange(ctx, deviceName, 0, -1).Result()
	if err != nil {
		klog.V(4).Infof("fail query data for deviceName,err:%v", err)
	}

	var dataModels []*common.DataModel

	for _, jsonStr := range dataJSON {
		var data common.DataModel
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			klog.V(4).Infof("Error unMarshaling data: %v\n", err)
			continue
		}

		dataModels = append(dataModels, &data)
	}
	return dataModels, nil

	//TODO implement me
	//panic("implement me")
}

func (d *DataBaseConfig) GetPropertyDataByDeviceName(deviceName string, propertyData string) ([]*common.DataModel, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DataBaseConfig) GetDataByTimeRange(start int64, end int64) ([]*common.DataModel, error) {
	ctx := context.Background()
	dataJSON, err := RedisCli.ZRangeByScore(ctx, "device2", &redis.ZRangeBy{
		Min: strconv.Itoa(int(start)),
		Max: strconv.Itoa(int(end)),
	}).Result()

	if err != nil {
		klog.V(4).Infof("fail query data: %v\n", err)
		return nil, err
	}

	var dataModels []*common.DataModel

	for _, jsonStr := range dataJSON {
		var data common.DataModel
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			klog.V(4).Infof("Error unMarshaling data: %v\n", err)
			continue
		}

		dataModels = append(dataModels, &data)
	}

	return dataModels, nil
	//TODO implement me
	//panic("implement me")
}

func (d *DataBaseConfig) DeleteDataByTimeRange(start int64, end int64) ([]*common.DataModel, error) {
	//TODO implement me
	panic("implement me")
}
