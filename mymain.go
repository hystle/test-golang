package main

import (
	"fmt"
	// "example/zach/testgolang/mystring"
	// "github.com/google/go-cmp/cmp"
	// "example/zach/testgolang/myyaml"
	// "example/zach/testgolang/mykube"
	// "example/zach/testgolang/mygoroutine"
	"example/zach/testgolang/mygit"
)

func main() {
	// fmt.Println(testStringPkg.ReverseRunes("Hello World"))
	// fmt.Println(cmp.Diff("Hello World", "Hello Go"))

	// fmt.Println("----")

	// testYamlPkg.ReadYamlConfig("sampleCfg.yaml")

	// fmt.Println("----")

	// var kube testKubePkg.Kube
	// if ok := kube.Connect(); !ok {
	// 	return
	// }

	// testKubePkg.GetPods(&kube)	// get pods every 10 seconds

	// testKubePkg.StartKubeResWatch(&kube, "kube-system", "storage-provisioner", "pod")	// watch res
	// select {} 	// wait on kube watch

	// fmt.Println("----")

	// c := make(chan int)
	// quit := make(chan int)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		// print 10 numbers whenever available in chan c
	// 		fmt.Printf(" --> output: %d\n", <-c)
	// 	}
	// 	quit <- 0
	// }()
	// testGoroutinePkg.FibWithChan(c, quit)

	// fmt.Println("----")

	// fmt.Println(testGoroutinePkg.Same(testGoroutinePkg.New(3),testGoroutinePkg.New(3)))
	// fmt.Println(testGoroutinePkg.Same(testGoroutinePkg.New(3),testGoroutinePkg.New(5)))

	fmt.Println("----")
	// r, _ := testGit.GitCloneRepo("/Users/ziczhou/Desktop/test", "https://github.com/go-git/go-git.git")
	// testGit.GitCheckoutBranch(r, "master")
	// testGit.GitCheckoutBranch(r, "wasm")
	// testGit.GitGetRemoteUrl(r, "/Users/ziczhou/Desktop/test")

	r, _ := testGit.GitOpenRepo("/Users/ziczhou/Desktop/test")
	err := testGit.GitPullRepo(r)
	if err == nil {
		fmt.Println("ok")
	}
}
