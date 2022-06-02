# gon.hcl
#
# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME

source = [ "./dist/macos_darwin_amd64_v1/godown" ]
bundle_id = "com.github.hirokistring.godown"

apple_id {
  # The environment variable "AC_USERNAME" must be defined with your Apple ID. 
  password = "@env:AC_PASSWORD"
}

sign {
  # The value must be changed to yours.
  application_identity = "Developer ID Application: Hiroki Saito"
}

zip {
  output_path = "./dist/godown_macos.zip"
}
