package aliyun

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/google/uuid"
)

// 配置获取
func GetRegionID() string {
	region := os.Getenv("ALIYUN_REGION_ID")
	if region == "" {
		return "cn-beijing"
	}
	return region
}

func GetSecurityGroupID() string {
	return os.Getenv("SECURITY_GROUP_ID")
}

func GetVSwitchID() string {
	return os.Getenv("VSWITCH_ID")
}

// CreateInstanceArgs 创建实例参数
type CreateInstanceArgs struct {
	InstanceType    string
	ImageID         string
	SecurityGroupID string
	InstanceName    string
	Bandwidth       int
	DiskSize        int
	Password        string
	RegionID        string
}

// 创建 ECS 实例
func CreateInstance(args CreateInstanceArgs) (*Instance, error) {
	client, err := NewECSClient()
	if err != nil {
		return nil, err
	}

	request := ecs.CreateCreateInstanceRequest()
	request.Scheme = "https"
	request.RegionId = args.RegionID
	if args.RegionID == "" {
		request.RegionId = GetRegionID()
	}

	request.InstanceType = args.InstanceType
	request.ImageId = args.ImageID
	request.InstanceName = args.InstanceName

	if args.SecurityGroupID != "" {
		request.SecurityGroupId = args.SecurityGroupID
	} else {
		request.SecurityGroupId = GetSecurityGroupID()
	}

	request.VSwitchId = GetVSwitchID()

	// 计费方式：按量付费
	request.InstanceChargeType = "PostPaid"
	request.SystemDisk.Category = "cloud_essd"

	response, err := client.CreateInstance(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	instanceID := response.InstanceId

	// 启动实例
	startRequest := ecs.CreateStartInstanceRequest()
	startRequest.InstanceId = instanceID
	_, err = client.StartInstance(startRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to start instance: %w", err)
	}

	return &Instance{
		ID:           instanceID,
		Name:         args.InstanceName,
		Status:       "Running",
		InstanceType: args.InstanceType,
		CreatedTime:  time.Now(),
	}, nil
}

// 删除实例
func DeleteInstance(instanceID string) error {
	client, err := NewECSClient()
	if err != nil {
		return err
	}

	request := ecs.CreateDeleteInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID
	request.ForceStop = requests.NewBoolean(true)

	_, err = client.DeleteInstance(request)
	return err
}

// 列出所有实例
func ListAllInstances() ([]Instance, error) {
	client, err := NewECSClient()
	if err != nil {
		return nil, err
	}
	return client.ListInstances()
}

// 获取单个实例
func GetOneInstance(instanceID string) (*Instance, error) {
	client, err := NewECSClient()
	if err != nil {
		return nil, err
	}
	return client.GetInstance(instanceID)
}

// 启动实例
func StartOneInstance(instanceID string) error {
	client, err := NewECSClient()
	if err != nil {
		return err
	}
	return client.StartInstance(instanceID)
}

// 停止实例
func StopOneInstance(instanceID string) error {
	client, err := NewECSClient()
	if err != nil {
		return err
	}
	return client.StopInstance(instanceID)
}

// 获取监控数据
func GetOneInstanceMetrics(instanceID string) (map[string]interface{}, error) {
	client, err := NewECSClient()
	if err != nil {
		return nil, err
	}
	return client.GetInstanceMetrics(instanceID)
}

// 获取实例规格选项
func GetInstanceTypeOptions() []InstanceTypeOption {
	return []InstanceTypeOption{
		{Code: "ecs.n1.small", Name: "1核1G", CPU: 1, Memory: 1024, Disk: 40},
		{Code: "ecs.n1.medium", Name: "1核2G", CPU: 1, Memory: 2048, Disk: 40},
		{Code: "ecs.n2.small", Name: "2核2G", CPU: 2, Memory: 2048, Disk: 40},
		{Code: "ecs.n2.medium", Name: "2核4G", CPU: 2, Memory: 4096, Disk: 40},
		{Code: "ecs.n4.small", Name: "2核4G", CPU: 2, Memory: 4096, Disk: 40},
		{Code: "ecs.n4.medium", Name: "2核8G", CPU: 2, Memory: 8192, Disk: 40},
		{Code: "ecs.n4.large", Name: "4核8G", CPU: 4, Memory: 8192, Disk: 40},
		{Code: "ecs.n4.xlarge", Name: "4核16G", CPU: 4, Memory: 16384, Disk: 40},
	}
}

// InstanceTypeOption 实例规格选项
type InstanceTypeOption struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	CPU     int    `json:"cpu"`
	Memory  int    `json:"memory"`
	Disk    int    `json:"disk"`
}

// 实例信息转 JSON
func (i *Instance) ToJSON() string {
	data, _ := json.Marshal(i)
	return string(data)
}

// ============================================================================
// 以下是原有的 ECS Client 代码
// ============================================================================

// 阿里云 ECS 客户端
type ECSClient struct {
	client       *ecs.Client
	regionID     string
	securityGroup string
}

// 实例信息
type Instance struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Status       string            `json:"status"`
	CPU          int               `json:"cpu"`
	Memory       int               `json:"memory"` // MB
	PublicIP     string            `json:"public_ip"`
	PrivateIP    string            `json:"private_ip"`
	InstanceType string            `json:"instance_type"`
	CreatedTime  time.Time         `json:"created_time"`
	Tags         map[string]string `json:"tags"`
}

// 创建 ECS 客户端
func NewECSClient() (*ECSClient, error) {
	accessKeyID := os.Getenv("ALIYUN_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
	regionID := os.Getenv("ALIYUN_REGION_ID")
	if regionID == "" {
		regionID = "cn-beijing"
	}

	if accessKeyID == "" || accessKeySecret == "" {
		return nil, fmt.Errorf("ALIYUN_ACCESS_KEY_ID or ALIYUN_ACCESS_KEY_SECRET not set")
	}

	client, err := ecs.NewClientWithAccessKey(regionID, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create ECS client: %w", err)
	}

	return &ECSClient{
		client:       client,
		regionID:     regionID,
		securityGroup: os.Getenv("SECURITY_GROUP_ID"),
	}, nil
}

// 创建 ECS 实例
func (c *ECSClient) CreateInstance(name, instanceType string, cpu, memory int) (*Instance, error) {
	// 请求创建实例
	request := ecs.CreateCreateInstanceRequest()
	request.Scheme = "https"
	request.RegionId = c.regionID
	request.InstanceType = instanceType
	request.ImageId = "aliyun_2_1903_x64_20G_alibase_20231227.vhd" // 默认镜像
	request.InstanceName = name
	request.SecurityGroupId = c.securityGroup
	request.VSwitchId = os.Getenv("VSWITCH_ID")

	// 计费方式：按量付费
	request.InstanceChargeType = "PostPaid"
	request.SystemDisk.Category = "cloud_essd"

	// 安全组
	if c.securityGroup != "" {
		request.SecurityGroupId = c.securityGroup
	}

	response, err := c.client.CreateInstance(request)
	if err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	instanceID := response.InstanceId

	// 启动实例
	startRequest := ecs.CreateStartInstanceRequest()
	startRequest.InstanceId = instanceID
	_, err = c.client.StartInstance(startRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to start instance: %w", err)
	}

	return &Instance{
		ID:           instanceID,
		Name:         name,
		Status:       "Running",
		CPU:          cpu,
		Memory:       memory,
		InstanceType: instanceType,
		CreatedTime:  time.Now(),
	}, nil
}

// 查询实例列表
func (c *ECSClient) ListInstances() ([]Instance, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.RegionId = c.regionID
	request.PageSize = requests.NewInteger(100)

	response, err := c.client.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}

	instances := make([]Instance, 0, len(response.Instances.Instance))
	for _, inst := range response.Instances.Instance {
		// 跳过已释放的实例
		if inst.Status == "Deleted" || inst.Status == "Releasing" {
			continue
		}

		var publicIP, privateIP string
		if len(inst.PublicIpAddress.IpAddress) > 0 {
			publicIP = inst.PublicIpAddress.IpAddress[0]
		}
		if len(inst.InnerIpAddress.IpAddress) > 0 {
			privateIP = inst.InnerIpAddress.IpAddress[0]
		}

		createdTime, _ := time.Parse("2006-01-02T15:04:05Z", inst.CreationTime)

		// 解析标签
		tags := make(map[string]string)
		if inst.Tag != nil {
			for _, tag := range inst.Tag.Tag {
				tags[tag.TagKey] = tag.TagValue
			}
		}

		instances = append(instances, Instance{
			ID:           inst.InstanceId,
			Name:         inst.InstanceName,
			Status:       inst.Status,
			CPU:          inst.Cpu,
			Memory:       inst.Memory,
			PublicIP:     publicIP,
			PrivateIP:    privateIP,
			InstanceType: inst.InstanceType,
			CreatedTime:  createdTime,
			Tags:         tags,
		})
	}

	return instances, nil
}

// 查询单个实例
func (c *ECSClient) GetInstance(instanceID string) (*Instance, error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.RegionId = c.regionID
	request.InstanceIds = fmt.Sprintf(`["%s"]`, instanceID)

	response, err := c.client.DescribeInstances(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	if len(response.Instances.Instance) == 0 {
		return nil, fmt.Errorf("instance not found")
	}

	inst := response.Instances.Instance[0]
	var publicIP, privateIP string
	if len(inst.PublicIpAddress.IpAddress) > 0 {
		publicIP = inst.PublicIpAddress.IpAddress[0]
	}
	if len(inst.InnerIpAddress.IpAddress) > 0 {
		privateIP = inst.InnerIpAddress.IpAddress[0]
	}

	createdTime, _ := time.Parse("2006-01-02T15:04:05Z", inst.CreationTime)

	return &Instance{
		ID:           inst.InstanceId,
		Name:         inst.InstanceName,
		Status:       inst.Status,
		CPU:          inst.Cpu,
		Memory:       inst.Memory,
		PublicIP:     publicIP,
		PrivateIP:    privateIP,
		InstanceType: inst.InstanceType,
		CreatedTime:  createdTime,
	}, nil
}

// 启动实例
func (c *ECSClient) StartInstance(instanceID string) error {
	request := ecs.CreateStartInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID

	_, err := c.client.StartInstance(request)
	return err
}

// 停止实例
func (c *ECSClient) StopInstance(instanceID string) error {
	request := ecs.CreateStopInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID

	_, err := c.client.StopInstance(request)
	return err
}

// 重启实例
func (c *ECSClient) RebootInstance(instanceID string) error {
	request := ecs.CreateRebootInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID

	_, err := c.client.RebootInstance(request)
	return err
}

// 删除实例
func (c *ECSClient) DeleteInstance(instanceID string) error {
	request := ecs.CreateDeleteInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID
	request.ForceStop = requests.NewBoolean(true)

	_, err := c.client.DeleteInstance(request)
	return err
}

// 查询实例监控数据
func (c *ECSClient) GetInstanceMetrics(instanceID string) (map[string]interface{}, error) {
	// 阿里云云监控 API
	request := ecs.CreateDescribeInstanceMonitorDataRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID
	request.StartTime = requests.NewFormattedInt64(time.Now().Add(-1 * time.Hour).Unix() * 1000)
	request.EndTime = requests.NewFormattedInt64(time.Now().Unix() * 1000)

	response, err := c.client.DescribeInstanceMonitorData(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	metrics := make(map[string]interface{})
	if len(response.MonitorData.MonitorData) > 0 {
		data := response.MonitorData.MonitorData[0]
		metrics["cpu"] = data.CPU
		metrics["memory"] = data.Memory
		metrics["disk_read"] = data.DiskReadBps
		metrics["disk_write"] = data.DiskWriteBps
		metrics["network_in"] = data.NetRx
		metrics["network_out"] = data.NetTx
	}

	return metrics, nil
}

// 创建快照
func (c *ECSClient) CreateSnapshot(instanceID, diskID string) (string, error) {
	request := ecs.CreateCreateSnapshotRequest()
	request.Scheme = "https"
	request.DiskId = diskID
	request.SnapshotName = fmt.Sprintf("snapshot-%s-%s", instanceID, uuid.New().String()[:8])

	response, err := c.client.CreateSnapshot(request)
	if err != nil {
		return "", fmt.Errorf("failed to create snapshot: %w", err)
	}

	return response.SnapshotId, nil
}

// 模拟客户端（当没有真实凭证时使用）
type MockECSClient struct{}

func NewMockECSClient() *MockECSClient {
	return &MockECSClient{}
}

func (c *MockECSClient) CreateInstance(name, instanceType string, cpu, memory int) (*Instance, error) {
	return &Instance{
		ID:           "i-" + uuid.New().String()[:8],
		Name:         name,
		Status:       "Running",
		CPU:          cpu,
		Memory:       memory,
		InstanceType: instanceType,
		PublicIP:     fmt.Sprintf("47.92.%d.%d", uuid.New().Fields()[0].(uint32)%256, uuid.New().Fields()[1].(uint32)%256),
		CreatedTime:  time.Now(),
	}, nil
}

func (c *MockECSClient) ListInstances() ([]Instance, error) {
	return []Instance{
		{
			ID:           "i-mock-001",
			Name:         "prod-api",
			Status:       "Running",
			CPU:          2,
			Memory:       4096,
			PublicIP:     "47.92.100.101",
			PrivateIP:    "172.16.0.101",
			InstanceType: "ecs.n4.small",
			CreatedTime:  time.Now().Add(-24 * time.Hour),
		},
		{
			ID:           "i-mock-002",
			Name:         "prod-db",
			Status:       "Running",
			CPU:          4,
			Memory:       8192,
			PublicIP:     "47.92.100.102",
			PrivateIP:    "172.16.0.102",
			InstanceType: "ecs.n4.large",
			CreatedTime:  time.Now().Add(-48 * time.Hour),
		},
	}, nil
}

func (c *MockECSClient) GetInstance(instanceID string) (*Instance, error) {
	instances, _ := c.ListInstances()
	for _, inst := range instances {
		if inst.ID == instanceID {
			return &inst, nil
		}
	}
	return nil, fmt.Errorf("instance not found")
}

func (c *MockECSClient) StartInstance(instanceID string) error {
	return nil
}

func (c *MockECSClient) StopInstance(instanceID string) error {
	return nil
}

func (c *MockECSClient) RebootInstance(instanceID string) error {
	return nil
}

func (c *MockECSClient) DeleteInstance(instanceID string) error {
	return nil
}

func (c *MockECSClient) GetInstanceMetrics(instanceID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"cpu":    45.5,
		"memory": 62.3,
		"disk":   38.2,
		"net_in":  125.5,
		"net_out": 85.3,
	}, nil
}

// 工厂函数：优先尝试真实客户端，失败时使用模拟
func NewECSClientOrMock() interface {
	CreateInstance(name, instanceType string, cpu, memory int) (*Instance, error)
	ListInstances() ([]Instance, error)
	GetInstance(instanceID string) (*Instance, error)
	StartInstance(instanceID string) error
	StopInstance(instanceID string) error
	RebootInstance(instanceID string) error
	DeleteInstance(instanceID string) error
	GetInstanceMetrics(instanceID string) (map[string]interface{}, error)
} {
	client, err := NewECSClient()
	if err != nil {
		fmt.Printf("[Aliyun] Using mock client: %v\n", err)
		return NewMockECSClient()
	}
	fmt.Println("[Aliyun] Using real ECS client")
	return client
}

// 实例信息转 JSON
func (i *Instance) ToJSON() string {
	data, _ := json.Marshal(i)
	return string(data)
}
