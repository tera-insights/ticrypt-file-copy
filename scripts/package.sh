#!/bin/bash
# This script builds an RPM package of the TICP

version=`scripts/getVersion.sh`
package="ticp"

spec=package/ticp.spec
tarFile=~/rpmbuild/SOURCES/${package}-v${version}.tar
directory="${package}-${version}/"

echo "Packaging"
tar -cvf ${tarFile} bin/ --transform "s,^,${directory},"
tar -rvf ${tarFile} package/ --transform "s,^,${directory},"
# tar -rvf ${tarFile} doc/ --transform "s,^,${directory},"

echo "Compressing"
gzip -f ${tarFile}

echo "Building"
rpmbuild -vv -ba $spec