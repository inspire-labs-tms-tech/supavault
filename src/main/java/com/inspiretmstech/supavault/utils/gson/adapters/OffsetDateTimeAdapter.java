package com.inspiretmstech.supavault.utils.gson.adapters;

import com.google.gson.TypeAdapter;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;

import java.io.IOException;
import java.time.OffsetDateTime;
import java.time.format.DateTimeFormatter;

public class OffsetDateTimeAdapter extends TypeAdapter<OffsetDateTime> {
    private static final DateTimeFormatter formatter = DateTimeFormatter.ISO_OFFSET_DATE_TIME;

    @Override
    public void write(JsonWriter out, OffsetDateTime value) throws IOException {
        out.value(value != null ? value.format(formatter) : null);
    }

    @Override
    public OffsetDateTime read(JsonReader in) throws IOException {
        String dateTimeString = in.nextString();
        return dateTimeString != null ? OffsetDateTime.parse(dateTimeString, formatter) : null;
    }
}
