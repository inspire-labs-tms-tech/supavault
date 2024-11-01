package com.inspiretmstech.supavault.utils.gson;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.inspiretmstech.supavault.utils.gson.adapters.OffsetDateTimeAdapter;

import java.time.OffsetDateTime;

public class GSON {

    public static final Gson GLOBAL = new GsonBuilder()
            .registerTypeAdapter(OffsetDateTime.class, new OffsetDateTimeAdapter())
            .create();

}
