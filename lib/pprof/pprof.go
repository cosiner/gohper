// Copy from http://1234n.com/?post/wgskfs
package pprof

// case "lookup heap":
//           p := pprof.Lookup("heap")
//           p.WriteTo(os.Stdout, 2)
//       case "lookup threadcreate":
//           p := pprof.Lookup("threadcreate")
//           p.WriteTo(os.Stdout, 2)
//       case "lookup block":
//           p := pprof.Lookup("block")
//           p.WriteTo(os.Stdout, 2)
//       case "start cpuprof":
//           if cpuProfile == nil {
//               if f, err := os.Create("game_server.cpuprof"); err != nil {
//                   log.Printf("start cpu profile failed: %v", err)
//               } else {
//                   log.Print("start cpu profile")
//                   pprof.StartCPUProfile(f)
//                   cpuProfile = f
//               }
//           }
//       case "stop cpuprof":
//           if cpuProfile != nil {
//               pprof.StopCPUProfile()
//               cpuProfile.Close()
//               cpuProfile = nil
//               log.Print("stop cpu profile")
//           }
//       case "get memprof":
//           if f, err := os.Create("game_server.memprof"); err != nil {
//               log.Printf("record memory profile failed: %v", err)
//           } else {
//               runtime.GC()
//               pprof.WriteHeapProfile(f)
//               f.Close()
//               log.Print("record memory profile")
//           }
