# Why was trivrost developed?
Despite an increasingly attractive web development landscape, established desktop applications (and those which are actually well-suited to being one) continue to provide business value in lieu of the potential risks and gains of a more modern remake. However, getting an application onto a user's computer and then keeping it up to date (all while being secure) is more difficult than one would hope for it to be: unconventional system configurations and strict company policies will call for you to take action. This is to be expected especially if you do not know some of your users directly; for example when they are partners of your client. trivrost makes sure that an always-online application is always up to date, can be deployed on a computer with one click, and has [numerous battle-tested solutions](troubleshooting.md) ready for you for when support requests come in. The project was further motivated by the need for a replacement for the now decommissioned *Java Web Start*.

# Why "trivrost"?
Originally, we wanted to call the project "Bifröst", which, in nordic mythology,  is a "rainbow bridge" connecting the great beyond (your webserver) with the earthly realm (the users' computers). We changed that up a little to prevent the name from containing an umlaut and clashing with the names of other projects, yielding "trivrost".

# About Java Web Start
Java Web Start (JWS) was a technology to deploy Java applications. A user would download a `.jnlp` file, which would be processed by the JWS component of the user's existing Java installation as soon as they double-clicked it, causing required files to be downloaded and the application to run. Unfortunately, JWS is discontinued starting with Java SE 11; even if it wasn't, it would still be limited to Java applications. Moreover, there no longer are Java Runtime Environment installations available for end users and Java's update mechanism has been removed along with its control panel. Instead of using a system-wide Java installation, Java applications should now ship with their desired Java runtime. For more details see:
* https://www.oracle.com/technetwork/java/javase/javaclientroadmapupdate2018mar-4414431.pdf
* https://docs.oracle.com/en/java/javase/11/migrate/index.html

To bundle a Java application with a JRE/JDK, there are several alternatives:
* JLink/JPackager: These tools can produce a small JVM containing only the modules you need. However, if you can't modularize your application in a good way, that JVM will be quite large. Also, JPackager was removed in Java 11.
* [Graal VM](https://www.graalvm.org/docs/reference-manual/aot-compilation/): Graal VM is a really cool project. Beside other stuff, it contains the possibility to compile your Java application to native code. That subproject is called Native Images. The problem is that currently (November 2018), Native Images is only supported on Linux and MacOS.
* [Launch4J](http://launch4j.sourceforge.net/): This project only produces `.exe` files, so you can't deploy on Linux and MacOS.
* [getdown](https://github.com/threerings/getdown/): A cool tool to replace JWS. It updates the app and can bundle a JRE/JDK. This project was the main inspiration source of trivrost, along with JWS itself. However, getdown itself is written in Java, so it relies on an installed JRE/JDK on the user's system.

Despite the project being motivated by a Java technology, note that trivrost can download and launch any executable.
