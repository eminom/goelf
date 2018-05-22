package compress

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"

	log "../logs"
)

func loadfrom(file string, tw *tar.Writer) error {
	fin, ef := os.Open(file)
	if ef != nil {
		return ef
	}
	defer fin.Close()
	_, err := io.Copy(tw, fin)
	return err
}

func CreateTar(files []string, tr func(string) string, writer io.Writer) error {
	gw := gzip.NewWriter(writer)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		log.Printf("%v", file)

		info, ei := os.Stat(file)
		if ei != nil {
			return ei
		}
		// symlink should be ""
		header, eh := tar.FileInfoHeader(info, "")
		if eh != nil {
			return eh
		}
		header.Name = tr(file)
		if ew := tw.WriteHeader(header); ew != nil {
			return ew
		}
		if el := loadfrom(file, tw); el != nil {
			return el
		}
	}
	return nil
}

// outd: output directory
func Detar(input io.Reader, outd string) error {
	gr, err := gzip.NewReader(input)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	//defer tr.Close() // no Close for tar.Reader
	for {
		header, err := tr.Next()
		if err != nil {
			if io.EOF == err {
				break
			}
			return err
		}
		outF := path.Join(outd, header.Name)
		log.Printf("%v", outF)
		baseD := path.Dir(outF)
		err = os.MkdirAll(baseD, os.ModePerm)
		if err != nil {
			return err
		}
		fout, err := os.Create(outF)
		if err != nil {
			return err
		}
		io.Copy(fout, tr)
		err = fout.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
