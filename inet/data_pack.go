package inet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"gmr/gmr-server/config"
	"gmr/gmr-server/iface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度
func (dp *DataPack) GetHeadLen() uint32 {
	//dataLen uint32(4字节) + id uint32(4字节)
	return 8
}

//封包方法，先写长度、类型，再写消息内容
//dataLen|id|data
func (dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	//写入数据缓冲区
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将id写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	//将data写入dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法,先读到包长度和包类型，再读取消息内容
func (dp *DataPack) Unpack(binaryData []byte) (iface.IMessage, error) {
	//读入数据缓冲区
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	//将dataLen读入dataBuff
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//将id读入dataBuff
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//将data读入dataBuff
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}
	//判断datalen是否超出最大包长
	if config.GetPackageSize() > 0 && msg.DataLen > config.GetPackageSize() {
		return nil, errors.New("msg size to large")
	}

	return msg, nil
}
