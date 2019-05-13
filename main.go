package main

import (
  "fmt"
  "time"
  "os/exec"
  "log"

  //"github.com/anjannath/systray"
  "../systray_sm_fork_pim"

  "hugo-control/assets"
  "hugo-control/config"
  "hugo-control/hugo"
  "github.com/kr/pretty"
)

var (
  siteSubmenus = make(map[string]*systray.MenuItem)
  menuItemLiveUrl *systray.MenuItem
  menuItemToggleHugoServer *systray.MenuItem
  menuItemOpenConcept *systray.MenuItem
  menuItemExit *systray.MenuItem
  menuSelectSite *systray.MenuItem
)

func main() {
  config.SetCurrentSite()
  log.Printf("Eugene Config: %# v", pretty.Formatter(config.CurrentSite))
  systray.Run(onReady, onExit)
}

func setCurrentSiteMenu(){

  //open Live Url
  if (config.CurrentSite.Live_Url != "") {
   menuItemLiveUrl = systray.AddMenuItem(fmt.Sprintf("Open %s", config.CurrentSite.Live_Url ) , "", 0)
  }

  //start Server Concept
  //stop Server Concept
  if hugo.HugoRunning(){
    menuItemToggleHugoServer = systray.AddMenuItem("stop lokale server", "", 0)
  } else {
    menuItemToggleHugoServer = systray.AddMenuItem("start lokale server", "", 0)
  }

  //open concept versie
  menuItemOpenConcept = systray.AddMenuItem(fmt.Sprintf("Open %s in conceptversie", config.CurrentSite.Name), "", 0)

}

func switchSitesMenu(){
  if(len(config.CurrentConfig.Sites)>1){
    menuSelectSite = systray.AddSubMenu("Switch site")
    for _, site := range config.CurrentConfig.Sites {
      tmpSubmenuItem := menuSelectSite.AddSubMenuItem(site.Name,"", 0)
      siteSubmenus[site.Name] = tmpSubmenuItem
    }

    for sName, sMenuItem := range siteSubmenus {
      go func(name string, siteMenuitem *systray.MenuItem) {
        for {
          <-siteMenuitem.OnClickCh()
          log.Println("Selecting %s", name)
          config.FindSiteIndexByName(name)
          config.SetCurrentSiteIndexByName(name)
          if hugo.HugoRunning(){
            hugo.KillHugo();
          }
          updateSiteMenu()
        }
      }(sName, sMenuItem)
    }
  }
}

func updateSiteMenu() {
  //open Live Url
  if (config.CurrentSite.Live_Url != "") {
   menuItemLiveUrl.SetTitle(fmt.Sprintf("Open %s", config.CurrentSite.Live_Url ))
  }

  //open concept versie
  menuItemOpenConcept.SetTitle(fmt.Sprintf("Open %s in conceptversie", config.CurrentSite.Name))
}

func renderMenu(){

  setCurrentSiteMenu()
  systray.AddSeparator()
  switchSitesMenu()
  systray.AddSeparator()
  menuItemExit = systray.AddMenuItem("Quit", "", 0)

  listenToServer()
  handleMenuClicks()
}

func listenToServer(){
  go func() {
    for {
      time.Sleep(time.Second)
      if hugo.HugoRunning(){
        menuItemToggleHugoServer.SetTitle("Stop Server")
        menuItemOpenConcept.Enable()
      } else{
        menuItemToggleHugoServer.SetTitle("Start Server")
        menuItemOpenConcept.Disable()
      }
    }
  }()
}

func handleMenuClicks(){
  go func() {
    for {
      select {

      case <-menuItemLiveUrl.OnClickCh():
        exec.Command("/usr/bin/open", config.CurrentSite.Live_Url).Output()
      case <-menuItemOpenConcept.OnClickCh():
        exec.Command("/usr/bin/open", "http://localhost:1313").Output()
      case <-menuItemToggleHugoServer.OnClickCh():
        if hugo.HugoRunning(){
          hugo.KillHugo();
        } else{
          hugo.StartHugo();
          menuItemOpenConcept.Show()
        }
      case <-menuItemExit.OnClickCh():
        systray.Quit()
        return
      }
    }
  }()
}

func onReady() {
  fmt.Printf("OnReady: %v+\n", time.Now())
  systray.SetIcon(images.MonoData)

  if(config.FatalError != "") {
    systray.AddMenuItem(config.FatalError,"",0)
  } else {
    renderMenu()
  }
}

func onExit() {
  if hugo.HugoRunning(){
    hugo.KillHugo();
  }
}
