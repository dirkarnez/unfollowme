unfollowme
==========
```javascript
const oldArray = [{"id": "1"}, {"id": "2"}];
const newArray = [{"id": "1"}];

oldArray.filter(oldElement => !newArray.map(newElement => newElement["id"]).includes(oldElement["id"]))
```
