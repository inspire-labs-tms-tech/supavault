package com.inspiretmstech.supavault.bases;

@FunctionalInterface
public interface Executor<T> {

     T execute() throws Exception;

}
