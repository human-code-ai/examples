package com.humancodeai.springbootjava.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class VerifyResult {
    @JsonProperty("human_id")
    private String humanId;
}