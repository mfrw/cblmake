# cblmake
A simple PoC based on CBL-Mariner


### How to build from source:
```bash
git clone https://github.com/mfrw/cblmake
cd cblmake
git checkout dev

cd cmd/cblmake
go build -v -tags netgo
```

### How to install using go get:

```bash
go install github.com/mfrw/cblmake/cmd/cblmake@dev
```
# TODO:

- [ ] Embed the worker chroot
- [ ] Remove hard-coded values
- [x] Finish `packSrpms`
- [x] Finish `extractSRPMS`
- [x] Finish `createGraph`
- [x] Finish `resolvePackages`
- [ ] Finish `buildPackages`
- [ ] Finish `buildImage`
- [ ] Finish `buildISO`
