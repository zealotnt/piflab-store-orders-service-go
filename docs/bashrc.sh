# Utility for API container
# Copy to ~/.bashrc
# Run `source ~/.bashrc` to apply
# or just: source docs/bashrc.sh

ApiRun(){
gin -p 80 run
}
gU(){
goose up
}
gD(){
goose down
}
gR(){
goose redo
}
