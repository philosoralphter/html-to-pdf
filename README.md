# html-to-pdf
##Status
Curently works but deplo9yemnt is broken.  when in docker, chrome fails with error:

`Fontconfig warning: "/etc/fonts/fonts.conf", line 146: blank doesn't take any effect anymore. please remove it from your fonts.conf`

Then it makes an illegible output.  This was fixed inthe past but the commit was lost when a machine was exchanged.

##Past Issues
###fonts config <blank> issue
Some issue with the fonts config file that comes with centos.  Workaround is to instlal some other fonts package.  Hve elected to install the google noto fonts. see https://bugs.chromium.org/p/chromium/issues/detail?id=695212


