# Terminal cmd

Cho phép chạy cmd theo dạng tree.

## Mô tả

Khi có quá nhiều command phức tạp cần chạy, alias có một số hạn chế như phải nhớ lệnh. tcmd cung cấp chức năng tương tự, kèm theo một giao diện dạng tree. Có thể group cmd theo tính năng vào một folder.

## Format

Create file/folder trong folder cmd với định dạng: 

```go
{
  "cmd": "echo",
  "inputs": ["Hello"]
}
```

Trong đó cmd là command muốn chạy, inputs là các inputs cho command đó.
