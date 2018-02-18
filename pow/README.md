# Proof of work exercise

Write two methods, one called `mint` and one called `verify`.

* `mint`
  - Mint should take two arguments: a `challenge` (random integer) and a `work_factor` (number of leading 0s in the hash).
  - It should return a `token`, which is a random string such that SHA2(`challenge` || `token`) starts with at least `work_factor` many 0s.
  - Use hex encoding rather than binary encoding for simplicity. (You'll want no more than 4 for your work factor.)
* `verify`
  - This should take three arguments: the `challenge`, the `work_factor`, and the `token`.
  - It should return `true` or `false` based on whether the token is valid.

Bonus: if you have extra time, add timestamping and implement a cache of recent tokens so that that double-spends are rejected.
