/**
 * Copyright (c) 2025 by SAP Labs Bulgaria,
 * url: http://www.sap.com
 * All rights reserved.
 * This software is the confidential and proprietary information
 * of SAP SE, Walldorf. You shall not disclose such Confidential
 * Information and shall use it only in accordance with the terms
 * of the license agreement you entered into with SAP.
 * Created on Nov 12, 2025 by C5407138
 */
package common

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpPort int
}

func DotEnv() {
	if err := godotenv.Load(".env"); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}
}
