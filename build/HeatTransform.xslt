<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform" xmlns:wix="http://schemas.microsoft.com/wix/2006/wi"
                version="1.0"
                xmlns="http://schemas.microsoft.com/wix/2006/wi"
                exclude-result-prefixes="wix">

    <xsl:output method="xml" encoding="UTF-8" indent="yes"/>

    <xsl:template match="wix:Wix">
        <xsl:copy>
            <xsl:apply-templates select="@*"/>
            <xsl:apply-templates/>
        </xsl:copy>
    </xsl:template>

    <!-- ### Adding the Win64-attribute to all Components -->
    <xsl:template match="wix:Component">

        <xsl:copy>
            <xsl:apply-templates select="@*"/>
            <!-- Adding the Win64-attribute as we have a x64 application -->
            <xsl:attribute name="Win64">yes</xsl:attribute>

            <!-- Now take the rest of the inner tag -->
            <xsl:apply-templates select="node()"/>
        </xsl:copy>

    </xsl:template>

    <xsl:template match="@*|node()">
        <xsl:copy>
            <xsl:apply-templates select="@*|node()"/>
        </xsl:copy>
    </xsl:template>

</xsl:stylesheet>
