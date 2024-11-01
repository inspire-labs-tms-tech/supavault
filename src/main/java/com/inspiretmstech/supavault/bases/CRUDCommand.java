package com.inspiretmstech.supavault.bases;

import java.util.UUID;

public interface CRUDCommand {

    int list(boolean json);

    int delete(UUID id);

}
