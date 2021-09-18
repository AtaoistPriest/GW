package net

//#include <sys/ioctl.h>
//#include <net/if.h>
//#include <string.h>
//#include <unistd.h>
//#include <linux/ethtool.h>
//#include <linux/types.h>
//#include <linux/sockios.h>
//int SetEthtoolValue(unsigned char * dev, int cmd, int value)
//{
//	struct ifreq ifr;
//	int fd;
//	struct ethtool_value ethv;
//	fd = socket(AF_INET, SOCK_DGRAM, 0);
//	if (fd == -1)
//	{
//		return -1;
//	}
//	strncpy(ifr.ifr_name, dev, sizeof(ifr.ifr_name));
//	ethv.cmd = cmd;
//	ethv.data = value;
//	ifr.ifr_data = (void *) &ethv;
//	if (ioctl(fd, SIOCETHTOOL, (char *)&ifr) < 0)
//	{
//		close(fd);
//		return -1;
//	}
//	close(fd);
//	return 0;
//}
import "C"

func EthInit(ifaceStr string)  {
	ifaceName := []byte(ifaceStr)
	//Set RX hw csum disable
	C.SetEthtoolValue((*C.uchar)(&ifaceName[0]), C.int(21), C.int(0))
	//Set TX hw csum disable
	C.SetEthtoolValue((*C.uchar)(&ifaceName[0]), C.int(23), C.int(0))
	//Set GRO disable
	C.SetEthtoolValue((*C.uchar)(&ifaceName[0]), C.int(44), C.int(0))
	//Set TSO disable
	C.SetEthtoolValue((*C.uchar)(&ifaceName[0]), C.int(31), C.int(0))
}
