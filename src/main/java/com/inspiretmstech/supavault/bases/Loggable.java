package com.inspiretmstech.supavault.bases;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public abstract class Loggable {

    protected final Logger logger;

    public Loggable(Class<?> cls) {
        this.logger = LoggerFactory.getLogger(cls);
    }

}
