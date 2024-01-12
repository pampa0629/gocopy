package main
import (
	"io"
	"os"
	"flag"
	"fmt"
	"path/filepath"
)

var g_isRemove bool = false
var g_isInfo bool = false
var g_isHelp bool = false

func main() {
	sourceDir := flag.String("s", "", "source dir")
	destDir := flag.String("d", "", "dest dir")
	flag.BoolVar(&g_isRemove, "r", false, "is remove source")
	flag.BoolVar(&g_isInfo, "i", true, "is output info")
	flag.BoolVar(&g_isHelp, "h", false, "is show help")
    flag.Parse()

	if g_isHelp {
		flag.Usage()
	} else if *sourceDir == "" || *destDir == "" {
		fmt.Println("source dir or dest dir is empty")
	} else if err := copyDir(*sourceDir, *destDir, 0); err != nil {
		panic(err)
	}
}
 
func copyDir(source, dest string, tab int) error {
	strTab := ""
	for i:=0; i<tab; i++ {
		strTab += "  "
	}

	// 创建目标文件夹
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	// 获取源文件夹中所有的文件和子文件夹
	files, err := os.ReadDir(source)
	if err != nil {
		fmt.Print(fmt.Sprintf("ReadDir: %s error.\n Error is: %s\n", source, err))
		return err
	}	
	// 遍历所有文件和子文件夹
	for _, file := range files {
		src := filepath.Join(source, file.Name())
		dst := filepath.Join(dest, file.Name())
		// 判断类型，如果是文件夹或者软连接，则递归复制
		// fileInfo,_ := os.Lstat(file)
		if file.Type()&os.ModeSymlink != 0 {
			if err := copySymlink(src, dst); err != nil {
				fmt.Print(fmt.Sprintf("copySymlink: %s to %s error.\n Error is: %s\n", src, dst, err))
				return err
			}
		} else if file.Type()&os.ModeDir != 0 {
			if g_isInfo {
				fmt.Print(fmt.Sprintf("%sdir: %s \n", strTab, file.Name()))
			}
			if err := copyDir(src, dst, tab+1); err != nil {
				fmt.Print(fmt.Sprintf("copyDir: %s to %s error.\n Error is: %s\n", src, dst, err))
				return err
			}
		} else {
			// 复制文件
			if err := copyFile(src, dst); err != nil {
				fmt.Print(fmt.Sprintf("copyFile: %s to %s error.\n Error is: %s\n", src, dst, err))
				return err
			}
		}
	}

	if g_isRemove {
		err:=os.RemoveAll(source)
		if err != nil {
			fmt.Print(fmt.Sprintf("os.RemoveAll: %s error.\n Error is: %s\n", source, err))
			return err
		}
	}
	
	return nil
}

func copyFile(source string, dest string) error {
	in, err := os.Open(source)
	if err != nil {
		fmt.Print(fmt.Sprintf("Open: %s error.\n Error is: %s\n", source, err))
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		fmt.Print(fmt.Sprintf("Create: %s error.\n Error is: %s\n", source, err))
		return err
	}
	defer out.Close()
	if _, err = io.Copy(out, in); err != nil {
		fmt.Print(fmt.Sprintf("io Copy: %s to %s error.\n Error is: %s\n", in, out, err))
		return err
	}

	if g_isRemove {
		err:=os.Remove(source)
		if err != nil {
			fmt.Print(fmt.Sprintf("os.Remove: %s error.\n Error is: %s\n", source, err))
			return err
		}
	}
	return out.Close()
}

// 判断链接是否存在，存在的话就删掉
func removeSymlink(link string) error {
	_, err := os.Stat(link)
	if os.IsNotExist(err) {
		return nil
	} else {
		err:=os.Remove(link)
		return err
	}
}

func copySymlink(sourcePath, destinationPath string) error {
	// 获取源文件的符号链接目标路径
	linkTarget, err := os.Readlink(sourcePath)
	if err != nil {
		fmt.Print(fmt.Sprintf("Readlink: %s error.\n Error is: %s\n", sourcePath, err))
		return err
	}

	removeSymlink(destinationPath)
	// 创建新的符号链接
	err = os.Symlink(linkTarget, destinationPath)
	if err != nil {
		fmt.Print(fmt.Sprintf("Symlink: %s error.\n Error is: %s\n", destinationPath, err))
		return err
	}

	if g_isRemove {
		err:=os.Remove(sourcePath)
		if err != nil {
			fmt.Print(fmt.Sprintf("os.Remove: %s error.\n Error is: %s\n", sourcePath, err))
			return err
		}
	}

	return nil
}
