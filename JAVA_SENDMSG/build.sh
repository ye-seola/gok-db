cd "$(dirname "$0")"
javac SendMsg.java -cp android-30.jar && $ANDROID_HOME/build-tools/34.0.0/d8 *.class && mv classes.dex sendmsg.dex
rm -rf *.class