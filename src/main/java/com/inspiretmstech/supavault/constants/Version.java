package com.inspiretmstech.supavault.constants;

import picocli.CommandLine;

public class Version implements CommandLine.IVersionProvider {

    public static final String version = "0.0.0";

    @Override
    public String[] getVersion() throws Exception {
        return new String[]{version};
    }

}
