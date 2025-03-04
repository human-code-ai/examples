package com.humancodeai.springbootjava.model;

import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class ClientResponse<T> {
    private int code;
    private String msg;
    private T result;
}