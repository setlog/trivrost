# Why was trivrost developed?
Despite an increasingly attractive web development landscape, desktop applications continue to exist and provide value for users and businesses. However, getting your application onto a user's computer and then keeping it up to date is a headache. This is especially true if you don't know some of your users directly; for example when they are partners of your client. Even before then, unexpectable and obscure system configurations **will** cause your application to fail to start. trivrost was motivated by the need for a replacement for the now decommissioned *Java Web Start*.

Java Web Start (JWS) was a technology to deploy Java applications. A user would download a `.jnlp` file, which would be processed by the JWS component of the user's existing Java installation as soon as they double-clicked it, causing required files to be downloaded and the application to run. Unfortunately, JWS is discontinued starting with Java SE 11; and even if it wasn't, it would be limited to Java applications. Moreover, there no longer are Java Runtime Environment installations available for end users and Java's update mechanism has been removed along with its control panel. Instead of using a system-wide Java installation, Java applications should now ship with their desired Java runtime. For more details see:
* https://www.oracle.com/technetwork/java/javase/javaclientroadmapupdate2018mar-4414431.pdf
* https://docs.oracle.com/en/java/javase/11/migrate/index.html

To bundle a Java application with a JRE/JDK, there are several alternatives:
* JLink/JPackager: These tools can produce a small JVM containing only the modules you need. However, if you can't modularize your application in a good way, that JVM will be quite large. Also, JPackager was removed in Java 11.
* [Graal VM](https://www.graalvm.org/docs/reference-manual/aot-compilation/): Graal VM is a really cool project. Beside other stuff, it contains the possibility to compile your Java application to native code. That subproject is called Native Images. The problem is that currently (November 2018), Native Images is only supported on Linux and MacOS.
* [Launch4J](http://launch4j.sourceforge.net/): This project only produces `.exe` files, so you can't deploy on Linux and MacOS.
* [getdown](https://github.com/threerings/getdown/): A cool tool to replace JWS. It updates the app and can bundle a JRE/JDK. This project was the main inspiration source of trivrost, along with JWS itself. However, getdown itself is written in Java, so it relies on an installed JRE/JDK on the user's system.

Despite the project being motivated by a Java technology, note that trivrost can download and launch any executable.
