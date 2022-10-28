unfollowme
==========
### Notes
- use [dirkarnez/json-repl](https://github.com/dirkarnez/json-repl)
  ```javascript
  const oldArray = globalObj[1].content; // check devtools' console for new array index
  const newArray = globalObj[0].content; // check devtools' console for new array index
  console.table(oldArray.filter(oldElement => !newArray.map(newElement => newElement["id"]).includes(oldElement["id"])));
  ```
- error ```invalid character '<' looking for beginning of value```: make sure the account has not been kicked logout and update `seed.txt` file
