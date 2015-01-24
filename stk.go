package main

func main() {
	result, _ := Search("drush failed")
	for _, item := range result.Items {
		println(item.Title)
	}
}
