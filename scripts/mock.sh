for file in ./pkg/*
	do
		if [ "${file}" == "./pkg/ent" ]; then
			continue
		fi
		echo mockery --output ../pkg/mocks/ --recursive --all --dir ${file}
		mockery --output ../pkg/mocks/ --recursive --all --dir ${file}
	done