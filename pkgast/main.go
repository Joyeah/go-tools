// 通过加载并解析Go包中的公有变量
// 对于资源类的变量，批量生成对应的代码片段
/*
如提取gioui的icons包中的所有图标数据变量，生成想要的代码片段。
Item{"Action3DRotation", layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			icon, _ := widget.NewIcon(icons.Action3DRotation)
			return material.IconButton(th, &widget.Clickable{}, icon, "Action3DRotation").Layout(gtx)
		})},
*/
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
)

func main() {
	flag.Parse()
	// 配置加载包
	cfg := &packages.Config{Mode: packages.NeedSyntax | packages.NeedTypes}
	pkgs, err := packages.Load(cfg, "golang.org/x/exp/shiny/materialdesign/icons")
	if err != nil {
		log.Fatalf("加载包失败: %v", err)
	}

	// 遍历包
	for _, pkg := range pkgs {
		fmt.Printf("包名: %s\n", pkg.Name)
		
		// 遍历包中的语法树
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				decl, ok := n.(*ast.GenDecl)
				if !ok || decl.Tok != token.VAR {
					return true
				}

				for _, spec := range decl.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					// fmt.Println(valueSpec.Values)
					for _, name := range valueSpec.Names {
						if name.IsExported() {
							// fmt.Printf("公有变量: %s\n", name.Name)
							// fmt.Printf("icon, _ := widget.NewIcon(icons.%s)\n", name.Name)
							// fmt.Printf(tmpl, name.Name, name.Name, name.Name)
							fmt.Printf("IconItem{\"%s\", icons.%s, widget.Clickable{}},\n", name.Name, name.Name)

							if valueSpec.Type != nil {
								fmt.Printf("类型: %v\n", valueSpec.Type)
							}
						}
					}
				}
				return true
			})
		}
	}
}