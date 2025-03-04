package com.humancodeai.springbootjava.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class GetSessionIdResult {
    @JsonProperty("session_id")
    private String sessionId;
}
