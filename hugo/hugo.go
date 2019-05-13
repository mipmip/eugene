package hugo

import (
  "hugo-control/config"
  "os"
  "os/exec"
  "fmt"
  "path"
  "path/filepath"
)

var HugoServer = exec.Command(HugoBinPath(), "server","-D", "-s", HugoDir() )

func HugoBinExists() bool {

  hugofile := path.Clean(filepath.Join(config.Dir(), "hugo"))

  if stat, err := os.Stat(hugofile); err == nil && stat.Mode().IsRegular() {
    return true
  } else {
    return false
  }
}

func HugoDir() string {
  return config.CurrentSite.Hugo_Src_Dir
}

func HugoBinPath() string {
  return path.Clean(filepath.Join(config.Dir(), "hugo"))
}

func StartHugo() {
  HugoServer = exec.Command(HugoBinPath(), "server","-D", "-d",config.CurrentSite.Hugo_Output_Dir, "-s", HugoDir())
  HugoServer.Start()
}

func KillHugo()  {
  if(HugoServer.Process != nil){
    if err := HugoServer.Process.Kill(); err != nil {
      fmt.Println("failed to kill process: ")
    }
  } else {
    fmt.Println("trying alternative way to kill HUGO")
    exec.Command("bash", "-c", fmt.Sprintf("/bin/kill %s",HugoPid())).Output()
  }
}

func HugoPid() string {
  out, _ := exec.Command("bash", "-c", "/bin/ps ax | /usr/bin/grep \"bitbar-hugo\\/hugo\"| /usr/bin/grep -v grep | /usr/bin/head -n1 | /usr/bin/cut -d\" \" -f 1").Output()

  pid := string(out)
  if fmt.Sprintf("%s", pid) != "<nil>" {
    return pid
  } else {
    return ""
  }
}

func HugoRunning() bool {
  if HugoPid() != "" {
    return true
  } else {
    return false
  }
}




