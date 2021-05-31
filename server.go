package main

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"mace_convert/convert"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var convert_cmd = [3]string{"python", "/mace/tools/converter.py", "convert"}   // 命令

func init()  {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

}

//sentence := fmt.Sprintf("My Name is %s", name)
func main(){
	// 服务端转化 mace 模型
	register()

}


func register(){
	// 注册服务
	listen, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil{
		log.Fatalf("count not listen port: %v", "8080")
	}
	grpc_server := grpc.NewServer()
	convert.RegisterConvertServer(grpc_server, &mace{})
	err = grpc_server.Serve(listen)
	if err != nil {
		log.Fatalf("fail to start grpc server because %v", err)
	}
}
type mace_output map[string] string  // 命令返回的消息

type mace struct {
	// 实现服务的结构体 ，他里面需要实现全部的 Services  定义的 rpc 函数
	//按照生成的 .pb.go 文件接口定义的函数签名来写
	convert.UnimplementedConvertServer
}


func (cc *mace) Mace(ctx context.Context, rec *convert.Client) (*convert.Server, error){
	// mace 函数实现
	defer func(){
		err := recover()
		if err != nil{
			fmt.Print("mace 函数运行错误")
		}
	}()
	ymal_path := rec.GetPath() // 绝对路径
	dest_path := rec.GetDestPath()  // xxx.zip
	convert_cmd_new := make([]string, 3)
	_ = copy(convert_cmd_new, convert_cmd[:])
	stdout := run_cmd(append(convert_cmd_new, ymal_path)...)            // 运行命令
	model_save_path := extract_dest_dir(ymal_path)

	dir2zip(model_save_path, dest_path)   // 压缩文件
	return  send(stdout), nil
}


func run_cmd(args...string) mace_output{
	// 运行命令行
	// 执行命令并获得输出
	arg := []string{ }  // 初始化参数切片
	pri_mace_output := mace_output{"stdout": "", "stderr": ""}
	cmd_str := ""  // 默认参数
	var stdout, stderr bytes.Buffer   // 默认有初始值

	for i, v := range args{
		if i == 0{

			cmd_str = v
		}
		if i >0{
			if strings.HasSuffix(v, "yml"){
				arg = append(arg, "--config="+v)
				continue
			}
			arg = append(arg, v)
		}
	}
	cmd := exec.Command(cmd_str, arg...)  // 可变参数
	fmt.Print(cmd_str, arg)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil || len(stderr.String()) != 0{
		pri_mace_output["stderr"] = stderr.String()
		//log.Fatalf("命令运行发生严重错误   ： %v", stderr.String())  // Fatalf 会直接终止程序运行了
	}
	pri_mace_output["stdout"] = stdout.String()

	return pri_mace_output
}


func send(out mace_output) *convert.Server {
	// 发送消息
	return &convert.Server{Stdout: out["stdout"], Stderr: out["stderr"]}
}


func extract_dest_dir(path string) (abs_path string){
	file, err := os.Open(path) // For read access.
	if err != nil{
		fmt.Printf("打开文件 %v 错误 %v", path, err)
	}
	defer file.Close()
	// 读取整个文件
	file_byte_content, err := ioutil.ReadAll(file)
	if err != nil{
		log.Fatalf("读取文件的时候发生错误 %v", err)
	}
	file_content := string(file_byte_content)
	// 正则提取 关键字的文件夹内容
	regex := regexp.MustCompile(`library_name:(.*)`)  // library_name: mobilenet
	if regex == nil{
		log.Fatalf("REGX COMPILE ERROR")
	}

	if len(regex.FindStringSubmatch(file_content)) > 1{

		library_name := strings.Trim(regex.FindStringSubmatch(file_content)[1], " ")
		fmt.Println("解析的目标文件目录名字是:", library_name)
		abs_path = fmt.Sprintf("/mace/build/%v/model", library_name)
	}else{
		log.Fatalf("没找到保存的转换后的文件")
	}
	return
}