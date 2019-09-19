Why was trivrost developed?
-------------------------

Java Web Start (JWS) was a technology to deploy java applications. A user had to download a jnlp file, which contained information, so JWS could update the program every time a user tries to start it. Unfortunately JWS is not available under Java SE 11 and newer versions. Moreover, Java 11 shouldn't be installed on the end user systems. There are no JRE downloads for end users, there is no updater mechanism for Java anymore and the control panel was removed. Instead of using a system wide installed JRE and JWS, the applications should bundle a Java runtime. For more details see:
* https://www.oracle.com/technetwork/java/javase/javaclientroadmapupdate2018mar-4414431.pdf
* https://docs.oracle.com/en/java/javase/11/migrate/index.html

To bundle a Java application with a JRE/JDK, there are several alternatives:
* JLink/JPackager: These tools can produce a small JVM containing only the modules you need. However, if you can't modularize your application in a good way, that JVM will be quite large. Also, JPackager was removed in Java 11.
* [Graal VM](https://www.graalvm.org/docs/reference-manual/aot-compilation/): Graal VM is a really cool project. Beside other stuff, it contains the possibility to compile your Java application to native code. That subproject is called Native Images. The problem is that currently (November 2018), Native Images is only supported on Linux and MacOS.
* [Launch4J](http://launch4j.sourceforge.net/): This project only produces `.exe` files, so you can't deploy on Linux and MacOS.
* [getdown](https://github.com/threerings/getdown/): A cool tool to replace JWS. It updates the app and can bundle a JRE/JDK. This project was the main inspiration source of trivrost, along with JWS itself. However, getdown itself is written in Java, so it relies on an installed JRE/JDK on the user's system.

trivrost was developed to bundle an OpenJDK together with a Java application, to keep both – the OpenJDK and the Java application itself – up-to-date and to actually launch the application afterwards. However, it is designed so that it can be used in other scenarios as well.
